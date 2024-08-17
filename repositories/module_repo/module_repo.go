package modulerepo

import (
	"doko/gvn-ultimate-bot/models"

	"gorm.io/gorm"
)

type ModuleRepo interface {
	ListModules() ([]*models.AppModule, error)
	ActivateOrDisableModule(id uint, newStatus uint8) (*models.AppModule, error)
	CreateModule(name string, label string, activated uint8) (*models.AppModule, error)
	UpdateModuleConfig(id uint, newConfig string) (*models.AppModule, error)
}

type moduleRepo struct {
	db *gorm.DB
}

func NewAppModuleRepo(db *gorm.DB) ModuleRepo {
	return &moduleRepo{
		db: db,
	}
}

func (ar *moduleRepo) ListModules() ([]*models.AppModule, error) {
	var modules []*models.AppModule
	if err := ar.db.Find(&modules).Error; err != nil {
		return modules, err
	}

	return modules, nil
}

func (mr *moduleRepo) CreateModule(name string, label string, activated uint8) (*models.AppModule, error) {
	module := models.AppModule{IsActivated: activated, ModuleLabel: label, ModuleName: name}

	if err := mr.db.Create(&module).Error; err != nil {
		return &module, nil
	}

	return nil, nil
}

func (mr *moduleRepo) ActivateOrDisableModule(id uint, newStatus uint8) (*models.AppModule, error) {
	var module *models.AppModule
	if err := mr.db.Where("id = ?", id).First(&module).Error; err == nil {
		module.IsActivated = newStatus // 0 = disable, 1 = enable
		mr.db.Save(&module)
		return module, err
	}

	return nil, nil
}

func (mr *moduleRepo) UpdateModuleConfig(id uint, newConfig string) (*models.AppModule, error) {
	var module *models.AppModule
	if err := mr.db.Where("id = ?", id).First(&module).Error; err == nil {
		module.CustomConfig = newConfig // Must be valid JSON
		mr.db.Save(&module)
		return module, err
	}

	return nil, nil
}
