package gql

import (
	"doko/gvn-ultimate-bot/gql/gen"
	"doko/gvn-ultimate-bot/services/authservice"
	"doko/gvn-ultimate-bot/services/userservice"
)

type Resolver struct {
	UserService userservice.UserService
	AuthService authservice.AuthService
}

func (r *Resolver) Mutation() gen.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() gen.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
