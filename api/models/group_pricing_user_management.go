package models

type GPUser struct {
	ID              int64  `json:"id" gorm:"primary_key"`
	Email           string `json:"email"`
	LicenseId       string `json:"license_id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	FullName        string `json:"full_name"`
	RoleId          int    `json:"role_id"`
	EmailSignature  string `json:"email_signature" gorm:"type:text"`
}

type GPUserRole struct {
	ID          int            `json:"id" gorm:"primary_key"`
	RoleName    string         `json:"role_name" gorm:"unique"`
	Description string         `json:"description"`
	Permissions []GPPermission `json:"permissions" gorm:"many2many:role_permissions;"`
}

type GPPermission struct {
	ID             int    `json:"id" gorm:"primary_key"`
	Name           string `json:"name" gorm:"unique"`
	Slug           string `json:"slug" gorm:"unique"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	Tier           string `json:"tier"`
	ParentSlug     string `json:"parent_slug"`
	DisplayOrder   int    `json:"display_order"`
	Disabled       bool   `json:"disabled" gorm:"-"`
	DisabledReason string `json:"disabled_reason" gorm:"-"`
}
