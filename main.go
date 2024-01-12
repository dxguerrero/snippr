package main

import (
	"github.com/dxguerrero/snippr/contollers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	controllers.ReadFile()

	// Snippet routes
	router.GET("/snippet", controllers.GetSnippets)
	router.GET("/snippet/:id", controllers.GetSnippetByID)
	router.POST("/snippet", controllers.PostSnippet)

	// User routes
	router.GET("/user", controllers.GetUsers)
	router.POST("/user", controllers.PostUser)

	router.Run("localhost:8080")
}