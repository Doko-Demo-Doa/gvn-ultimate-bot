package models

import "gorm.io/gorm"

type AppModule struct {
	gorm.Model
	ModuleName   string `gorm:"NOT NULL;size:255"`
	ModuleLabel  string `gorm:"size:255"`
	IsActivated  uint8  `gorm:"DEFAULT:0"`
	CustomConfig string `gorm:"size:8192"`
}

func (AppModule) TableName() string {
	return "app_module"
}
