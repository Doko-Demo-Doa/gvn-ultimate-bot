package repositories

import (
	"doko/gin-sample/models"

	"gorm.io/gorm"
)

type Repo interface {
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
}

type userRepo struct {
	db *gorm.DB
}
