package main

import (
	"doko/gin-sample/app"
	"doko/gin-sample/configs"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := configs.GetConfig()
	println(config)

	app.Run()
}
