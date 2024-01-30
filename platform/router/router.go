package router


import (
	"encoding/gob"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/dxguerrero/snippr/platform/authenticator"
	"github.com/dxguerrero/snippr/platform/middleware"
	"github.com/dxguerrero/snippr/web/app/callback"
	"github.com/dxguerrero/snippr/web/app/login"
	"github.com/dxguerrero/snippr/web/app/logout"
	"github.com/dxguerrero/snippr/web/app/user"
	"github.com/dxguerrero/snippr/contollers" // Note to self: fix typo when right-click is working again.
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()
	controllers.ReadFile()
	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.Static("/public", "web/static")
	router.LoadHTMLGlob("web/template/*")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/user", user.Handler)
	router.GET("/logout", logout.Handler)

	snippetRoutes := router.Group("/snippet")
	snippetRoutes.Use(middleware.IsAuthenticated)
	
	// Snippet routes
	snippetRoutes.GET("/", controllers.GetSnippets)
	snippetRoutes.GET("/:id", controllers.GetSnippetByID)
	snippetRoutes.POST("/", controllers.PostSnippet)

	return router
}