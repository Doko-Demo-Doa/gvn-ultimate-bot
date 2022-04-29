package app

import (
	"doko/gin-sample/bot"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func Run() {
	bot.Bootstrap()
}
