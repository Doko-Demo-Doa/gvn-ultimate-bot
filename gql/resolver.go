package gql

import (
	"doko/gin-sample/services/authservice"
	"doko/gin-sample/services/userservice"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	UserService userservice.UserService
	AuthService authservice.AuthService
}

func (r *Resolver) Mutation() gen.MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct {
	*Resolver
}
