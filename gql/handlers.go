package gql

import (
	"doko/gvn-ultimate-bot/gql/gen"
	"doko/gvn-ultimate-bot/services/authservice"
	"doko/gvn-ultimate-bot/services/userservice"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func GraphqlHandler(us userservice.UserService, as authservice.AuthService) gin.HandlerFunc {
	conf := gen.Config{
		Resolvers: &Resolver{
			UserService: us,
			AuthService: as,
		},
	}

	exec := gen.NewExecutableSchema(conf)
	h := handler.New(exec)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func PlaygroundHandler(path string) gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", path)
	return func(c *gin.Context) { h.ServeHTTP(c.Writer, c.Request) }
}
