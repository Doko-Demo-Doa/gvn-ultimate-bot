package gql

import (
	"context"
	"doko/gvn-ultimate-bot/gql/gen"
	"errors"
)

func (r *queryResolver) User(ctx context.Context, id int) (*gen.User, error) {
	user, err := r.UserService.GetByID(uint(id))
	if err != nil {
		return nil, err
	}

	return &gen.User{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (r *queryResolver) UserProfile(ctx context.Context) (*gen.User, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, errors.New("unauthorized: Token is invalid")
	}

	user, err := r.UserService.GetByID(userID.(uint))
	if err != nil {
		return nil, err
	}

	return &gen.User{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
