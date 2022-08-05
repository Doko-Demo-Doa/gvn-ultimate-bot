package seeds

import (
	"doko/gin-sample/models"
	"doko/gin-sample/services/userservice"
	"doko/gin-sample/statics"

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
			Email:    "quan.pham@darenft.com",
			Password: string(hashed[:]),
			Role:     statics.AdminRole,
		},

		{
			Name:     "Giap Tran",
			Email:    "giap.tran@darenft.com",
			Password: string(hashed[:]),
			Role:     statics.Standard,
		},
	}

	for _, model := range users {
		us.Create(&model)
	}
}
