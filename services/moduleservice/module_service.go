package moduleservice

import (
	"doko/gvn-ultimate-bot/models"
	modulerepo "doko/gvn-ultimate-bot/repositories/modulerepo"
)

type ModuleService interface {
	ActivateOrDisableModule(id uint) (string, error)
	ListModules() ([]*models.AppModule, error)
}

type moduleService struct {
	Repo modulerepo.ModuleRepo
}

func NewModuleService(jwtSecret string) ModuleService {
	return &moduleService{}
}

func (ms *moduleService) ActivateOrDisableModule(id uint) (string, error) {
	panic("unimplemented")
}

func (ms *moduleService) ListModules() ([]*models.AppModule, error) {
	user, err := ms.Repo(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
