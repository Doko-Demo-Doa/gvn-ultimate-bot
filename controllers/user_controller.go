package controllers

import (
	"doko/gvn-ultimate-bot/models"
	"doko/gvn-ultimate-bot/services/authservice"
	"doko/gvn-ultimate-bot/services/userservice"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/resend/resend-go/v2"
)

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserOutput struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserUpdateInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserController interface {
	Register(*gin.Context)
	Login(*gin.Context)
	GetByID(*gin.Context)
	GetProfile(*gin.Context)
	Update(*gin.Context)
	ForgotPassword(*gin.Context)
	ResetPassword(*gin.Context)
}

type userController struct {
	us userservice.UserService
	as authservice.AuthService
}

func NewUserController(
	us userservice.UserService,
	as authservice.AuthService,
) UserController {
	return &userController{
		us: us,
		as: as,
	}
}

func (ctl *userController) ForgotPassword(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	resetToken, err := ctl.us.InitiateResetPassword(input.Email)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Send the email:
	resendApiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(resendApiKey)
	params := &resend.SendEmailRequest{
		From:    "Admin <info@aniviet.com>",
		To:      []string{input.Email},
		Text:    "Reset your password. Token is: " + resetToken,
		Subject: "Reset your password",
	}

	sent, err := client.Emails.Send(params)

	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	HTTPRes(c, http.StatusOK, "Email sent", sent.Id)
}

// GetByID implements UserController
func (ctl *userController) GetByID(c *gin.Context) {
	id, err := ctl.getUserID(c.Param(("id")))
	if err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user, err := ctl.us.GetByID(id)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			HTTPRes(c, http.StatusNotFound, es, nil)
			return
		}
		HTTPRes(c, http.StatusInternalServerError, es, nil)
		return
	}
	userOutput := ctl.mapToUserOutput(user)
	HTTPRes(c, http.StatusOK, "ok", userOutput)
}

// GetProfile implements UserController
func (ctl *userController) GetProfile(c *gin.Context) {
	id, exists := c.Get("user_id")
	if !exists {
		HTTPRes(c, http.StatusBadRequest, "Invalid User ID", nil)
		return
	}

	user, err := ctl.us.GetByID(id.(uint))
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	userOutput := ctl.mapToUserOutput(user)
	HTTPRes(c, http.StatusOK, "ok", userOutput)
}

// @Summary Login
// @Produce  json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/login [post]
func (ctrl *userController) Login(c *gin.Context) {
	var userInput UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	user, err := ctrl.us.GetByEmail(userInput.Email)

	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
	}

	err = ctrl.us.ComparePassword(userInput.Password, user.Password)

	if err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = ctrl.login(c, user)

	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
}

// Register implements UserController
func (ctl *userController) Register(c *gin.Context) {
	// Read user input
	var userInput UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	u := ctl.inputToUser(userInput)

	// Create user
	if err := ctl.us.Create(&u); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Login
	err := ctl.login(c, &u)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
}

func (ctl *userController) ResetPassword(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	token := c.Request.URL.Query().Get("token")
	if token == "" {
		HTTPRes(c, http.StatusNotFound, "Requires token", nil)
		return
	}

	user, err := ctl.us.CompleteUpdatePassword(token, input.Password)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	err = ctl.login(c, user)
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
}

// Update implements UserController
func (ctl *userController) Update(c *gin.Context) {
	// Get user id from context
	id, exists := c.Get("user_id")
	if !exists {
		HTTPRes(c, http.StatusBadRequest, "Invalid User ID", nil)
		return
	}

	// Retrieve user given id
	user, err := ctl.us.GetByID(id.(uint))
	if err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Read user input
	var userInput UserUpdateInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// Check user
	if user.ID != id {
		HTTPRes(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Update user record
	user.Name = userInput.Name
	user.Email = userInput.Email
	if err := ctl.us.Update(user); err != nil {
		HTTPRes(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Response
	userOutput := ctl.mapToUserOutput(user)
	HTTPRes(c, http.StatusOK, "ok", userOutput)
}

func (ctl *userController) getUserID(userIDParam string) (uint, error) {
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return 0, errors.New("user id should be a number")
	}
	return uint(userID), nil
}

func (ctl *userController) inputToUser(input UserInput) models.User {
	return models.User{
		Email:    input.Email,
		Password: input.Password,
	}
}

func (ctl *userController) mapToUserOutput(u *models.User) *UserOutput {
	return &UserOutput{
		ID:    u.ID,
		Email: u.Email,
		Role:  u.Role,
		Name:  u.Name,
	}
}

// Issue token and return user
func (ctl *userController) login(c *gin.Context, u *models.User) error {
	token, err := ctl.as.IssueToken(*u)
	if err != nil {
		return err
	}
	userOutput := ctl.mapToUserOutput(u)
	out := gin.H{"token": token, "user": userOutput}
	HTTPRes(c, http.StatusOK, "ok", out)
	return nil
}
