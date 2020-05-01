package routes

import (
	"rakoon/rakoon-back/authentication"
	"rakoon/rakoon-back/middleware"
	"rakoon/rakoon-back/utils"

	"github.com/gin-gonic/gin"
)

// InitRoutes calls the routes init
func InitRoutes(r *gin.Engine) {

	r.Use(middleware.CorsHandling)

	// Public routes
	public := r.Group("/v1")
	public.POST("/connect", func(c *gin.Context) { authentication.Connect(c) })
	public.POST("/subscribe", func(c *gin.Context) { authentication.Subscribe(c) })
	public.POST("/refresh/token", func(c *gin.Context) { authentication.RefreshToken(c) })

	// Private Routes
	private := r.Group("/v1")
	private.Use(middleware.JwtHandling)
	private.GET("/ping", func(c *gin.Context) { utils.Ping(c) })
	private.POST("/logout", func(c *gin.Context) { authentication.LogOut(c) })
}
