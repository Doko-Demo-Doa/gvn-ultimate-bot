package main

import (
	"doko/gin-sample/app"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app.Run(os.Getenv("DISCORD_TOKEN"))
}
