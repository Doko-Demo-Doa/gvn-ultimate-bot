package seeds

import (
	"doko/gvn-ultimate-bot/models"
	"doko/gvn-ultimate-bot/services/moduleservice"
	"doko/gvn-ultimate-bot/services/userservice"
	"doko/gvn-ultimate-bot/statics"

	"golang.org/x/crypto/bcrypt"
)

func SeedUsers(us userservice.UserService) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("testpwd"+"secret"), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	users := []models.User{
		{
			Name:     "Quan Pham",
			Email:    "example@example.com",
			Password: string(hashed[:]),
			Role:     statics.AdminRole,
		},
	}

	for _, model := range users {
		us.Create(&model)
	}
}

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
