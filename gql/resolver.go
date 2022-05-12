package gql

import (
	"doko/gin-sample/gql/gen"
	"doko/gin-sample/services/authservice"
	"doko/gin-sample/services/userservice"
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
