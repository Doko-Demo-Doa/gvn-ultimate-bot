package seeds

import (
	"doko/gvn-ultimate-bot/models"
	"doko/gvn-ultimate-bot/services/moduleservice"
)

func SeedModules(ms moduleservice.ModuleService) {
	modules := []models.AppModule{
		{
			ModuleName:  "pin_module",
			ModuleLabel: "Pin Module",
			IsActivated: 1,
		},
		{
			ModuleName:  "grant_role_module",
			ModuleLabel: "Grant Role Module",
			IsActivated: 1,
		},
	}

	for _, model := range modules {
		ms.CreateModule(model.ModuleName, model.ModuleLabel, model.IsActivated)
	}
}
