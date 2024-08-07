package app

import (
	"doko/gvn-ultimate-bot/bot"
	"doko/gvn-ultimate-bot/common/hmachash"
	"doko/gvn-ultimate-bot/common/randomstring"
	"doko/gvn-ultimate-bot/configs"
	"doko/gvn-ultimate-bot/controllers"
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	modulerepo "doko/gvn-ultimate-bot/repositories/module_repo"
	"doko/gvn-ultimate-bot/repositories/passwordreset"
	"doko/gvn-ultimate-bot/repositories/userrepo"

	"doko/gvn-ultimate-bot/seeds"
	"doko/gvn-ultimate-bot/services/authservice"
	"doko/gvn-ultimate-bot/services/discordservice"
	"doko/gvn-ultimate-bot/services/moduleservice"
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

	db.AutoMigrate(&models.User{}, &models.DiscordRole{}, &models.AppModule{})

	// Setup common
	rds := randomstring.NewRandomString()
	hm := hmachash.NewHMAC(config.HMACKey)

	// Setup repo
	userRepo := userrepo.NewUserRepo(db)
	moduleRepo := modulerepo.NewAppModuleRepo(db)

	pwdRepo := passwordreset.NewPasswordResetRepo(db)
	discordRepo := discordrepos.NewDiscordRoleRepo(db)

	// Setup services
	userService := userservice.NewUserService(userRepo, pwdRepo, rds, hm, config.Pepper)
	moduleService := moduleservice.NewModuleService(moduleRepo)
	authService := authservice.NewAuthService(config.JWTSecret)
	discordService := discordservice.NewDiscordService(discordRepo)

	// Seeding modules
	mModules, _ := moduleService.ListModules()
	if len(mModules) <= 0 {
		seeds.SeedModules(moduleService)
	}

	// Seeding users
	mUsers, _ := userService.ListUsers()
	if len(mUsers) <= 0 {
		seeds.SeedUsers(userService)
	}

	// Setup controllers
	userCtrl := controllers.NewUserController(userService, authService)
	moduleCtl := controllers.NewModuleController(moduleService)
	discordCtl := controllers.NewDiscordController(discordService)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

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
	discord.GET("/role/list", discordCtl.ListDiscordRoles)
	discord.POST("/role/create", discordCtl.CreateDiscordRole)

	// Module-related
	module := api.Group("/module")
	module.GET("/list", moduleCtl.ListModules)
	module.POST("/on-off", moduleCtl.ActivateOrDisableModule)

	port := fmt.Sprintf(":%s", config.Port)

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go func() {
		router.Run(port)
		wg.Done()
	}()

	go func() {
		// Bot setup
		bot.Bootstrap(db, discordService, moduleService)
		wg.Done()
	}()

	wg.Wait()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
