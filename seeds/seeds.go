package seeds

import (
	"doko/gin-sample/models"
	"doko/gin-sample/services/userservice"

	"golang.org/x/crypto/bcrypt"
)

func SeedUsers(us userservice.UserService) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("testpwd"), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	users := []models.User{
		{
			Name:     "Quan Pham",
			Email:    "quan.pham@darenft.com",
			Password: string(hashed[:]),
			Role:     "admin",
		},

		{
			Name:     "Giap Tran",
			Email:    "giap.tran@darenft.com",
			Password: string(hashed[:]),
			Role:     "admin",
		},
	}

	for _, model := range users {
		us.Create(&model)
	}
}
