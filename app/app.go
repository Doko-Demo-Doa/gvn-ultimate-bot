package app

import (
	"doko/gvn-ultimate-bot/bot"
	"doko/gvn-ultimate-bot/configs"
	"doko/gvn-ultimate-bot/controllers"
	"doko/gvn-ultimate-bot/models"
	discordrepos "doko/gvn-ultimate-bot/repositories/discord_repos"
	modulerepo "doko/gvn-ultimate-bot/repositories/module_repo"
	"doko/gvn-ultimate-bot/scheduler"
	"doko/gvn-ultimate-bot/seeds"
	"doko/gvn-ultimate-bot/services/adminaccessservice"
	"doko/gvn-ultimate-bot/services/discordservice"
	"doko/gvn-ultimate-bot/services/moduleservice"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
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

	// WARNING: Remember to run this the first time to create tables
	db.AutoMigrate(&models.DiscordRole{}, &models.DiscordRoleReactionEmbed{}, &models.AppModule{}, &models.DiscordUserRole{}, &models.AdminWhitelistedRole{})

	// Setup repo
	moduleRepo := modulerepo.NewAppModuleRepo(db)

	discordRepo := discordrepos.NewDiscordRoleRepo(db)
	discordRoleReactionEmbedRepo := discordrepos.NewDiscordRoleReactionEmbedRepo(db)
	discordUserRoleRepo := discordrepos.NewDiscordUserRoleRepo(db)
	adminWhitelistedRoleRepo := discordrepos.NewAdminWhitelistedRoleRepo(db)

	// Setup Discord state (shared between bot and scheduler)
	s := state.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentGuildMessages)
	s.AddIntents(gateway.IntentGuildMessageReactions)

	// Setup scheduler
	guildID := discord.GuildID(mustSnowflakeEnv("DISCORD_GUILD_ID"))

	// Setup services
	moduleService := moduleservice.NewModuleService(moduleRepo)
	discordRoleService := discordservice.NewDiscordRoleService(discordRepo, discordRoleReactionEmbedRepo, discordUserRoleRepo)
	discordRoleReactionEmbedService := discordservice.NewDiscordRoleReactionEmbedService(discordRoleReactionEmbedRepo, s, guildID)
	roleScheduler := scheduler.NewRoleScheduler(s, discordRoleService, guildID)
	adminAccessService := adminaccessservice.NewAdminAccessService(adminWhitelistedRoleRepo, s, guildID)

	// Seeding modules
	mModules, _ := moduleService.ListModules()
	if len(mModules) <= 0 {
		seeds.SeedModules(moduleService)
	}

	// Setup controllers
	moduleCtl := controllers.NewModuleController(moduleService)
	discordRoleCtl := controllers.NewDiscordController(discordRoleService, discordRoleReactionEmbedService, roleScheduler)
	adminAccessCtl := controllers.NewAdminAccessController(adminAccessService)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	// Setup routes
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ping")
	})

	api := router.Group("/api")

	// Discord-related APIs
	discord := api.Group("/discord")
	discord.GET("/channels", discordRoleCtl.ListDiscordChannels)
	discord.GET("/emojis", discordRoleCtl.ListDiscordEmojis)
	discord.GET("/role/list", discordRoleCtl.ListDiscordRoles)
	discord.POST("/role/create", discordRoleCtl.CreateDiscordRole)
	discord.POST("/role/assign", discordRoleCtl.AssignRoleToUser)
	discord.DELETE("/role/assign/:id", discordRoleCtl.RevokeRoleFromUser)
	discord.GET("/role/assignments", discordRoleCtl.ListRoleAssignments)

	discord.GET("/role-reaction/list", discordRoleCtl.ListDiscordRoleReactions)
	discord.GET("/role-reaction/:id", discordRoleCtl.GetDiscordRoleReaction)
	discord.POST("/role-reaction/upsert", discordRoleCtl.UpsertDiscordRoleReaction)
	discord.POST("/role-reaction/publish", discordRoleCtl.PublishDiscordRoleReaction)
	discord.DELETE("/role-reaction/:id", discordRoleCtl.DeleteDiscordRoleReaction)

	// Admin dashboard access control (Discord role whitelist)
	admin := api.Group("/admin")
	admin.GET("/access-check", adminAccessCtl.CheckAccess)
	admin.GET("/whitelisted-roles", adminAccessCtl.ListWhitelistedRoles)
	admin.POST("/whitelisted-roles", adminAccessCtl.UpsertWhitelistedRole)
	admin.DELETE("/whitelisted-roles/:id", adminAccessCtl.DeleteWhitelistedRole)

	// Module-related
	module := api.Group("/module")
	module.GET("/list", moduleCtl.ListModules)
	module.GET("/id/:id", moduleCtl.GetModuleByID)
	module.GET("/", moduleCtl.GetModuleByName)
	module.POST("/on-off", moduleCtl.ActivateOrDisableModule)
	module.POST("/update-config", moduleCtl.UpdateModuleConfig)

	port := fmt.Sprintf(":%s", config.Port)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		router.Run(port)
		wg.Done()
	}()

	go func() {
		bot.Bootstrap(s, discordRoleService, discordRoleReactionEmbedService, moduleService, roleScheduler)
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

func mustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(os.Getenv(env))
	if err != nil {
		panic(fmt.Sprintf("Invalid snowflake for $%s: %v", env, err))
	}
	return s
}
