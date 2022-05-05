// package main

// import (
// 	"doko/gin-sample/app"
// 	"log"

// 	"github.com/joho/godotenv"
// )

// func main() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	app.Run()
// }

package main

import (
	"github.com/umbracle/ethgo/jsonrpc"
)

func main() {
	client, err := jsonrpc.NewClient("https://mainnet.infura.io")
	if err != nil {
		panic(err)
	}

	client.Close()
}
