package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint
	Name     string `gorm:"size:255"`
	Email    string `gorm:"size:255"`
	Password string `gorm:"NOT NULL"`
	Role     string `gorm:"NOT NULL;size:255;DEFAULT:'standard'"`
}
