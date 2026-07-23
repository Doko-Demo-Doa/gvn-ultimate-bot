package models

import "gorm.io/gorm"

// AdminWhitelistedRole marks which Discord roles are permitted to access the admin dashboard.
type AdminWhitelistedRole struct {
	gorm.Model
	RoleNativeID string `gorm:"NOT NULL;unique;size:255;index"`
	AccessLevel  string `gorm:"NOT NULL;size:255;DEFAULT:'standard'"`
	Label        string `gorm:"size:255"`
}

func (AdminWhitelistedRole) TableName() string {
	return "admin_whitelisted_role"
}
