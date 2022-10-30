package userservice

import (
	"doko/gvn-ultimate-bot/common/hmachash"
	"doko/gvn-ultimate-bot/common/randomstring"
	"doko/gvn-ultimate-bot/models"
	pwdRepo "doko/gvn-ultimate-bot/repositories/passwordreset"
	"doko/gvn-ultimate-bot/repositories/userrepo"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(*models.User) error
	Update(*models.User) error
	HashPassword(rawPassword string) (string, error)
	ComparePassword(rawPassword string, passwordFromDB string) error
	InitiateResetPassowrd(email string) (string, error)
	CompleteUpdatePassword(token, newPassword string) (*models.User, error)
}

type userService struct {
	Repo    userrepo.Repo
	PwdRepo pwdRepo.Repo
	Rds     randomstring.RandomString
	hmac    hmachash.HMAC
	pepper  string
}

func (us *userService) ComparePassword(rawPassword string, passwordFromDB string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(passwordFromDB),
		[]byte(rawPassword+us.pepper),
	)
}

// CompleteUpdatePassword implements UserService
func (us *userService) CompleteUpdatePassword(token string, newPassword string) (*models.User, error) {
	hashedToken := us.hmac.Hash(token)

	pwr, err := us.PwdRepo.GetOneByToken(hashedToken)
	if err != nil {
		return nil, err
	}

	// If the password rest is over 1 hours old, it is invalid
	if time.Since(pwr.CreatedAt) > (1 * time.Hour) {
		return nil, errors.New("invalid token")
	}

	user, err := us.Repo.GetByID(pwr.UserID)
	if err != nil {
		return nil, err
	}

	hashedPass, err := us.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPass
	if err = us.Repo.Update(user); err != nil {
		return nil, err
	}

	if err = us.PwdRepo.Delete(pwr.ID); err != nil {
		fmt.Println("Failed to delete passwordreset record", pwr.ID, err.Error())
	}
	return user, nil
}

// Create implements UserService
func (us *userService) Create(user *models.User) error {
	hashedPass, err := us.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPass
	return us.Repo.Create(user)
}

func (us *userService) GetByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email(string) is required")
	}

	user, err := us.Repo.GetByEmail(email)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) GetByID(id uint) (*models.User, error) {
	if id == 0 {
		return nil, errors.New("id param is required")
	}

	user, err := us.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) HashPassword(rawPassword string) (string, error) {
	passAndPepper := rawPassword + us.pepper
	hashed, err := bcrypt.GenerateFromPassword([]byte(passAndPepper), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashed), err
}

func (us *userService) InitiateResetPassowrd(email string) (string, error) {
	user, err := us.Repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	token, err := us.Rds.GenerateToken()
	if err != nil {
		return "", err
	}

	hashedToken := us.hmac.Hash(token)

	pwd := models.PasswordReset{
		UserID: user.ID,
		Token:  hashedToken,
	}

	err = us.PwdRepo.Create(&pwd)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Update implements UserService
func (us *userService) Update(user *models.User) error {
	return us.Repo.Update(user)
}

func NewUserService(
	repo userrepo.Repo,
	pwdRepo pwdRepo.Repo,
	rds randomstring.RandomString,
	hmac hmachash.HMAC,
	pepper string,
) UserService {
	return &userService{
		Repo:    repo,
		PwdRepo: pwdRepo,
		Rds:     rds,
		hmac:    hmac,
		pepper:  pepper,
	}
}
