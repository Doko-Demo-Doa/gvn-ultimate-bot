package modulerepo

import (
	"doko/gvn-ultimate-bot/models"

	"gorm.io/gorm"
)

type ModuleRepo interface {
	ListModules() ([]*models.AppModule, error)
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

func (ar *moduleRepo) ActivateOrDisableModule(id uint) (*models.AppModule, error) {
	var module *models.AppModule
	if err := ar.db.Where("id = ?", id).First(&module).Error; err != nil {
		module.IsActivated = 1 // 0 = disable, 1 = enable
		ar.db.Save(&module)
		return module, err
	}

	return nil, nil
}
