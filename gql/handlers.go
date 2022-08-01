package gql

import (
	"doko/gin-sample/gql/gen"
	"doko/gin-sample/services/authservice"
	"doko/gin-sample/services/userservice"

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
