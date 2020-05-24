package routes

import (
	"rakoon/rakoon-back/controllers/authentication"
	"rakoon/rakoon-back/controllers/user"
	"rakoon/rakoon-back/controllers/utils"
	"rakoon/rakoon-back/middleware"

	"github.com/gin-gonic/gin"
)

// InitRoutes calls the routes init
func InitRoutes(r *gin.Engine) {

	r.Use(middleware.CorsHandling)

	// Public routes
	public := r.Group("/v1")

	public.POST("/subscribe", func(c *gin.Context) { user.Create(c) })
	public.POST("/connect", func(c *gin.Context) { user.Connect(c) })
	public.POST("/refresh/token", func(c *gin.Context) { authentication.RefreshToken(c) })

	// Private Routes
	private := r.Group("/v1")
	private.Use(middleware.JwtHandling)

	private.POST("/logout", func(c *gin.Context) { user.LogOut(c) })
	private.GET("/ping", func(c *gin.Context) { utils.Ping(c) })
	private.GET("/user/:id", func(c *gin.Context) { user.Get(c) })
	private.PUT("/user", func(c *gin.Context) { user.Update(c) })
	private.PUT("/user/password", func(c *gin.Context) { user.UpdatePassword(c) })
	private.PUT("/archive/user", func(c *gin.Context) { user.Archive(c) })
	private.DELETE("/user/:id", func(c *gin.Context) { user.Delete(c) })
}
