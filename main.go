package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DownloadEntity struct {
	id   string `json: "id"`
	link string `json: "link"`
	host string `json: "host"`
}

var sampleLinks = []DownloadEntity{
	{id: "1", link: "https://google.com", host: "rapidgator"},
}

func main() {
	router := gin.Default()
	router.GET("/downloads", getLinkById)
	router.POST("/albums", postLink)

	router.Run("localhost:8080")
}

func getLinks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, sampleLinks)
}

func postLink(c *gin.Context) {
	var newLink DownloadEntity

	// Bind the JSON to newLink
	if err := c.BindJSON(&newLink); err != nil {
		return
	}

	sampleLinks = append(sampleLinks, newLink)
	c.IndentedJSON(http.StatusCreated, newLink)
}

func getLinkById(c *gin.Context) {
	// id := c.Param("id")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "test message"})
}
