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
		{
			ModuleName:  "message_audit_module",
			ModuleLabel: "Message Audit Module",
			IsActivated: 0,
		},
		{
			ModuleName:  "member_sync_module",
			ModuleLabel: "Member Sync Module",
			IsActivated: 0,
		},
	}

	for _, model := range modules {
		existing, _ := ms.GetModuleByName(model.ModuleName)
		if existing == nil {
			ms.CreateModule(model.ModuleName, model.ModuleLabel, model.IsActivated)
		}
	}
}
