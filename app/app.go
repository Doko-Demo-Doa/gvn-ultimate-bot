package app

import (
	"doko/gin-sample/bot"
	"doko/gin-sample/configs"
	"doko/gin-sample/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	router = gin.Default()
)

func Run() {
	bot.Bootstrap()

	config := configs.GetConfig()

	db, err := gorm.Open(config.Postgres.GetPostgresConfigInfo())

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})
}
