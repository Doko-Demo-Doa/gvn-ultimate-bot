package models

import "gorm.io/gorm"

type PasswordReset struct {
	gorm.Model
	UserID uint   `gorm:"NOT NULL"`
	Token  string `gorm:"NOT NULL;unique_index"`
}
