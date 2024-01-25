package main

import (
	"github.com/dxguerrero/snippr/contollers"
	"github.com/dxguerrero/snippr/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	controllers.ReadFile()

	protectedRoutes := router.Group("/snippet")
	protectedRoutes.Use(middleware.AuthMiddleware())

	// Snippet routes
	protectedRoutes.GET("/", controllers.GetSnippets)
	protectedRoutes.GET("/:id", controllers.GetSnippetByID)
	protectedRoutes.POST("/", controllers.PostSnippet)

	// User routes
	userRoutes := router.Group("/user")
	userRoutes.GET("/", controllers.GetUser)
	userRoutes.POST("/", controllers.PostUser)
	userRoutes.POST("/login", controllers.Login)

	router.Run("localhost:8080")
}