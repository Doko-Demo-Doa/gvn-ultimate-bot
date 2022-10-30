package userrepo

import (
	"doko/gvn-ultimate-bot/models"

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

func NewUserRepo(db *gorm.DB) Repo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where(&models.User{Email: email}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) Update(user *models.User) error {
	return u.db.Save(user).Error
}
