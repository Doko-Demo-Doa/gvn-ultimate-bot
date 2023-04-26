package moduleservice

import (
	"doko/gvn-ultimate-bot/models"
	modulerepo "doko/gvn-ultimate-bot/repositories/module_repo"
)

type ModuleService interface {
	ActivateOrDisableModule(id uint, newStatus uint8) (*models.AppModule, error)
	ListModules() ([]*models.AppModule, error)
	CreateModule(name string, activated uint8) (*models.AppModule, error)
}

type moduleService struct {
	Repo modulerepo.ModuleRepo
}

func NewModuleService(repo modulerepo.ModuleRepo) ModuleService {
	return &moduleService{
		Repo: repo,
	}
}

func (ms *moduleService) ActivateOrDisableModule(id uint, newStatus uint8) (*models.AppModule, error) {
	module, err := ms.Repo.ActivateOrDisableModule(id, newStatus)
	if err != nil {
		return nil, err
	}

	return module, nil
}

func (ms *moduleService) CreateModule(name string, activated uint8) (*models.AppModule, error) {
	module, err := ms.Repo.CreateModule(name, activated)

	if err != nil {
		return nil, err
	}

	return module, nil
}

func (ms *moduleService) ListModules() ([]*models.AppModule, error) {
	user, err := ms.Repo.ListModules()
	if err != nil {
		return nil, err
	}

	return user, nil
}
