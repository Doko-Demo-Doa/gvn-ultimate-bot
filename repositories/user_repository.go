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

func NewUserRepo(db *gorm.DB) Repo {
	return &userRepo{
		db: db,
	}
}

// Create implements Repo
func (*userRepo) Create(user *models.User) error {
	panic("unimplemented")
}

// GetByEmail implements Repo
func (u *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID implements Repo
func (u *userRepo) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update implements Repo
func (*userRepo) Update(user *models.User) error {
	panic("unimplemented")
}
