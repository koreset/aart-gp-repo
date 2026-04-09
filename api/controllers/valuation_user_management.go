package controllers

import (
	"api/models"
	"api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetValPermissions(c *gin.Context) {
	permissions, err := services.GetValPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, permissions)
}

func GetValUserRoles(c *gin.Context) {
	roles, err := services.GetValUserRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, roles)
}

func CreateValUserRole(c *gin.Context) {
	var role models.ValUserRole
	err := c.BindJSON(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	role, err = services.CreateValUserRole(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, role)
}

func DeleteValUserRole(c *gin.Context) {
	roleId := c.Param("role_id")
	err := services.DeleteValUserRole(roleId)

	if err != nil && err.Error() == "role is in use" {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetValRolePermissions(c *gin.Context) {
	roleId := c.Param("role_id")
	permissions, err := services.GetValRolePermissions(roleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, permissions)
}

func AssignValRoleToUser(c *gin.Context) {
	var userRole models.OrgUser
	err := c.BindJSON(&userRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.AssignValRoleToUser(userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func RemoveValRoleFromUser(c *gin.Context) {
	var userRole models.OrgUser
	err := c.BindJSON(&userRole)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = services.RemoveValRoleFromUser(userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nil)
}

func GetValRoleForUser(c *gin.Context) {
	licenseId := c.Param("license_id")
	role, err := services.GetValRoleForUserLicense(licenseId)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, role)
}
