package gql

import (
	"context"
	"doko/gin-sample/gql/gen"
	"doko/gin-sample/models"
	"errors"
)

func (r *mutationResolver) Login(ctx context.Context, input gen.RegisterLogin) (*gen.RegisterLoginOutput, error) {
	user, err := r.UserService.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	err = r.UserService.ComparePassword(input.Password, user.Password)
	if err != nil {
		return nil, err
	}

	token, err := r.AuthService.IssueToken(*user)
	if err != nil {
		return nil, err
	}

	return &gen.RegisterLoginOutput{
		Token: token,
		User: &gen.User{
			ID:    int(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (r *mutationResolver) Register(ctx context.Context, input gen.RegisterLogin) (*gen.RegisterLoginOutput, error) {
	userDomain := &models.User{
		Email:    input.Email,
		Password: input.Password,
	}

	err := r.UserService.Create(userDomain)
	if err != nil {
		return nil, err
	}

	token, err := r.AuthService.IssueToken(*userDomain)
	if err != nil {
		return nil, err
	}

	return &gen.RegisterLoginOutput{
		Token: token,
		User: &gen.User{
			ID:    int(userDomain.ID),
			Name:  userDomain.Name,
			Email: userDomain.Email,
			Role:  userDomain.Role,
		},
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input gen.UpdateUser) (*gen.User, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, errors.New("unauthorized: Token is invlaid")
	}

	usr, err := r.UserService.GetByID(userID.(uint))
	if err != nil {
		return nil, err
	}

	if input.Email != "" {
		usr.Email = input.Email
	}

	if input.Name != nil {
		usr.Name = *input.Name
	}

	err = r.UserService.Update(usr)
	if err != nil {
		return nil, err
	}

	return &gen.User{
		ID:    int(usr.ID),
		Name:  usr.Name,
		Email: usr.Email,
		Role:  usr.Role,
	}, nil
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, errors.New("email is required")
	}

	// Issue token for user to update his/her password
	_, err := r.UserService.InitiateResetPassowrd(email)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ResetPassword(ctx context.Context, resetToken string, password string) (*gen.RegisterLoginOutput, error) {
	if resetToken == "" {
		return nil, errors.New("token is required")
	}

	if password == "" {
		return nil, errors.New("new password is required")
	}

	user, err := r.UserService.CompleteUpdatePassword(resetToken, password)
	if err != nil {
		return nil, err
	}

	token, err := r.AuthService.IssueToken(*user)
	if err != nil {
		return nil, err
	}

	return &gen.RegisterLoginOutput{
		Token: token,
		User: &gen.User{
			ID:    int(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
