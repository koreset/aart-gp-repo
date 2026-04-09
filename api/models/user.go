package models

import "time"

type AppUser struct {
	UserName  string `json:"userName" gorm:"type:varchar(255);uniqueIndex"`
	UserEmail string `json:"userEmail"`
}

type OrgUser struct {
	ID           int    `json:"id" gorm:"primary_key"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	LicenseId    string `json:"license_id"`
	Role         string `json:"role"`
	GPRole       string `json:"gp_role"`
	GPRoleId     int    `json:"gp_role_id"`
	ValRole      string `json:"val_role"`
	ValRoleId    int    `json:"val_role_id"`
	Organisation string `json:"organisation"`
}

type Organisation struct {
	Name string `json:"name"`
}

type Activity struct {
	ID          int       `json:"id" gorm:"primary_key"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserName    string    `json:"user_name"`
	UserEmail   string    `json:"user_email"`
	ObjectID    int       `json:"object_id"`
	ObjectType  string    `json:"object_type"`
	Date        time.Time `json:"date"`
}
