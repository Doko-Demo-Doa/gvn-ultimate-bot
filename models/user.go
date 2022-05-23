package models

import (
	"gorm.io/gorm"
)

// Field-level permissions: https://gorm.io/docs/models.html#Field-Level-Permission
type User struct {
	gorm.Model
	ID       uint
	Name     string `gorm:"size:255"`
	Email    string `gorm:"size:255"`
	Password string `gorm:"NOT NULL"`
	Role     string `gorm:"NOT NULL;size:255;DEFAULT:'standard'"`
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return "user"
}
