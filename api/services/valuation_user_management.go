package services

import (
	appLog "api/log"
	"api/models"
	"errors"
	"fmt"
	"strings"
)

// GetValPermissions returns all valuation permissions
func GetValPermissions() ([]models.ValPermission, error) {
	var permissions []models.ValPermission
	err := DB.Find(&permissions).Error
	if err != nil {
		return permissions, err
	}
	return permissions, nil
}

// GetValUserRoles returns all valuation roles with their permissions
func GetValUserRoles() ([]models.ValUserRole, error) {
	var roles []models.ValUserRole
	err := DB.Preload("Permissions").Find(&roles).Error
	if err != nil {
		return roles, err
	}
	return roles, nil
}

// CreateValUserRole creates or updates a valuation role with permissions
func CreateValUserRole(role models.ValUserRole) (models.ValUserRole, error) {
	permissions := role.Permissions
	role.Permissions = nil // avoid GORM trying to insert join table before Role exists

	if role.ID == 0 {
		if err := DB.Create(&role).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				return role, fmt.Errorf("role with name '%s' already exists", role.RoleName)
			}
			return role, err
		}
	} else {
		if err := DB.Save(&role).Error; err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				return role, fmt.Errorf("role with name '%s' already exists", role.RoleName)
			}
			return role, err
		}
	}

	// remove any previous associations
	if err := DB.Model(&role).Association("Permissions").Clear(); err != nil {
		return role, err
	}

	if err := DB.Model(&role).Association("Permissions").Append(permissions); err != nil {
		return role, err
	}

	return role, nil
}

// DeleteValUserRole deletes a valuation role if not in use
func DeleteValUserRole(roleId string) error {
	var role models.ValUserRole
	err := DB.Where("id = ?", roleId).First(&role).Error
	if err != nil {
		return err
	}

	// check if role is in use
	var users []models.OrgUser
	err = DB.Where("val_role_id = ?", roleId).Find(&users).Error
	if err != nil {
		fmt.Println(err)
	}
	if len(users) > 0 {
		return errors.New("role is in use")
	}

	// delete associated permissions
	err = DB.Model(&role).Association("Permissions").Clear()
	if err != nil {
		fmt.Println(err)
	}

	err = DB.Delete(&role).Error
	if err != nil {
		return err
	}
	return nil
}

// GetValRolePermissions returns permissions for a specific valuation role
func GetValRolePermissions(roleId string) ([]models.ValPermission, error) {
	var role models.ValUserRole
	err := DB.Where("id = ?", roleId).Preload("Permissions").First(&role).Error
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

// AssignValRoleToUser assigns a valuation role to a user
func AssignValRoleToUser(user models.OrgUser) error {
	err := DB.Save(&user).Error
	if err != nil {
		appLog.Error(err.Error())
		return err
	}
	return nil
}

// RemoveValRoleFromUser removes the valuation role from a user
func RemoveValRoleFromUser(user models.OrgUser) error {
	user.ValRoleId = 0
	user.ValRole = "None"
	err := DB.Where("id = ?", user.ID).Updates(&user).Error
	if err != nil {
		appLog.Error(err.Error())
		return err
	}
	return nil
}

// GetValRoleForUserLicense returns the valuation role for a user by license ID
func GetValRoleForUserLicense(licenseId string) (models.ValUserRole, error) {
	var orgUser models.OrgUser
	err := DB.Where("license_id = ?", licenseId).First(&orgUser).Error
	if err != nil {
		appLog.Error("Error getting user role for license: ", err.Error())
		return models.ValUserRole{}, err
	}
	var role models.ValUserRole
	err = DB.Where("id = ?", orgUser.ValRoleId).Preload("Permissions").First(&role).Error
	if err != nil {
		appLog.Error("Error getting valuation user role: ", err.Error())
		return models.ValUserRole{}, err
	}
	return role, nil
}

// MigrateValuationUserTables migrates valuation user management tables.
// Permissions and roles are seeded from JSON files via BaseData() in install.go.
func MigrateValuationUserTables() error {
	err := DB.AutoMigrate(
		&models.ValUserRole{},
		&models.ValPermission{},
	)
	if err != nil {
		return err
	}
	// Also ensure OrgUser has the new val_role fields
	return DB.AutoMigrate(&models.OrgUser{})
}
