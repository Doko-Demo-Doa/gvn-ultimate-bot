package app

import (
	"doko/gin-sample/bot"
	"doko/gin-sample/common/hmachash"
	"doko/gin-sample/common/randomstring"
	"doko/gin-sample/configs"
	"doko/gin-sample/models"
	"doko/gin-sample/repositories/userrepo"

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

	router.Run()

	config := configs.GetConfig()

	db, err := gorm.Open(config.Postgres.GetPostgresConfigInfo())

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	// Setup common
	rds := randomstring.NewRandomString()
	hm := hmachash.NewHMAC(config.HMACKey)

	userRepo := userrepo.NewUserRepo(db)

	// Setup services
	print(rds, hm, userRepo)

	// Bot setup
	bot.Bootstrap()
}
