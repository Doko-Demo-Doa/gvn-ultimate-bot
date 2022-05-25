package app

import (
	"doko/gin-sample/bot"
	"doko/gin-sample/common/hmachash"
	"doko/gin-sample/common/randomstring"
	"doko/gin-sample/configs"
	"doko/gin-sample/controllers"
	"doko/gin-sample/models"
	"doko/gin-sample/repositories/passwordreset"
	"doko/gin-sample/repositories/userrepo"
	"doko/gin-sample/services/authservice"
	"doko/gin-sample/services/userservice"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	router = gin.Default()
)

func Run() {
	/*
		====== Swagger setup ============
		(http://localhost:3000/swagger/index.html)
	*/
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	config := configs.GetConfig()

	db, err := gorm.Open(config.Postgres.GetPostgresConfigInfo())

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	// Setup common
	rds := randomstring.NewRandomString()
	hm := hmachash.NewHMAC(config.HMACKey)

	// Setup repo
	userRepo := userrepo.NewUserRepo(db)
	pwdRepo := passwordreset.NewPasswordResetRepo(db)

	println(config.Pepper)

	// Setup services
	userService := userservice.NewUserService(userRepo, pwdRepo, rds, hm, config.Pepper)
	authService := authservice.NewAuthService(config.JWTSecret)

	// Seeding
	// seeds.SeedUsers(userService)

	// Setup controllers
	userCtrl := controllers.NewUserController(userService, authService)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup routes
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ping")
	})

	// router.GET("/graphql", gql.Pl)

	api := router.Group("/api")

	api.POST("/register", userCtrl.Register)
	api.POST("/login", userCtrl.Login)
	api.POST("/forgot-password", userCtrl.ForgotPassword)
	api.POST("/reset-password", userCtrl.ResetPassword)

	user := api.Group("/users")

	user.GET("/:id", userCtrl.GetByID)

	account := api.Group("/account")
	// account.Use(middlewares)
	account.GET("/profile", userCtrl.GetProfile)
	account.PUT("/profile", userCtrl.Update)

	port := fmt.Sprintf(":%s", config.Port)
	router.Run(port)

	// Bot setup
	bot.Bootstrap()
}
