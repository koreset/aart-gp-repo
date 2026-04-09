package models

type ValUserRole struct {
	ID          int             `json:"id" gorm:"primary_key"`
	RoleName    string          `json:"role_name" gorm:"unique"`
	Description string          `json:"description"`
	Permissions []ValPermission `json:"permissions" gorm:"many2many:val_role_permissions;"`
}

type ValPermission struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Name        string `json:"name" gorm:"unique"`
	Slug        string `json:"slug" gorm:"unique"`
	Description string `json:"description"`
}
