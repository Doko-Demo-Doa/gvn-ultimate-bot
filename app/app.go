package app

import (
	"doko/gvn-ultimate-bot/bot"
	"doko/gvn-ultimate-bot/common/hmachash"
	"doko/gvn-ultimate-bot/common/randomstring"
	"doko/gvn-ultimate-bot/configs"
	"doko/gvn-ultimate-bot/controllers"
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	"doko/gvn-ultimate-bot/repositories/passwordreset"
	"doko/gvn-ultimate-bot/repositories/userrepo"
	"doko/gvn-ultimate-bot/services/authservice"
	"doko/gvn-ultimate-bot/services/discordservice"
	"doko/gvn-ultimate-bot/services/userservice"
	"fmt"
	"net/http"
	"sync"

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

	db, err := gorm.Open(config.Postgres.GetPostgresConfigInfo(), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{}, &models.DiscordRole{})

	// Setup common
	rds := randomstring.NewRandomString()
	hm := hmachash.NewHMAC(config.HMACKey)

	// Setup repo
	userRepo := userrepo.NewUserRepo(db)
	pwdRepo := passwordreset.NewPasswordResetRepo(db)
	discordRepo := discordrepos.NewDiscordRoleRepo(db)

	// Setup services
	userService := userservice.NewUserService(userRepo, pwdRepo, rds, hm, config.Pepper)
	authService := authservice.NewAuthService(config.JWTSecret)
	discordService := discordservice.NewDiscordService(discordRepo)

	// Seeding
	// seeds.SeedUsers(userService)

	// Setup controllers
	userCtrl := controllers.NewUserController(userService, authService)
	discordCtl := controllers.NewDiscordController(discordService)

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

	// Discord-related APIs
	discord := api.Group("/discord")
	discord.GET("/list-roles", discordCtl.ListDiscordRoles)
	discord.POST("/create-role", discordCtl.CreateDiscordRole)

	port := fmt.Sprintf(":%s", config.Port)

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go func() {
		router.Run(port)
		wg.Done()
	}()

	go func() {
		// Bot setup
		bot.Bootstrap(db, discordService)
		wg.Done()
	}()

	wg.Wait()
}
