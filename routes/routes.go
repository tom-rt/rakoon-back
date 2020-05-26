package routes

import (
	"rakoon/rakoon-back/controllers/authentication"
	"rakoon/rakoon-back/controllers/user"
	"rakoon/rakoon-back/middleware"

	"github.com/gin-gonic/gin"
)

// InitRoutes calls the routes init
func InitRoutes(r *gin.Engine) {

	r.Use(middleware.CorsHandling)

	// Public routes
	public := r.Group("/v1")

	public.POST("/user", func(c *gin.Context) { user.Create(c) })
	public.POST("/user/connect", func(c *gin.Context) { user.Connect(c) })
	public.POST("/refresh/token", func(c *gin.Context) { authentication.RefreshToken(c) })

	// Private Routes
	private := r.Group("/v1")
	private.Use(middleware.JwtHandling)

	private.GET("/user/:id", func(c *gin.Context) { user.Get(c) })
	private.PUT("/user/:id", func(c *gin.Context) { user.Update(c) })
	private.PUT("/user/:id/password", func(c *gin.Context) { user.UpdatePassword(c) })
	private.PUT("/user/:id/logout", func(c *gin.Context) { user.LogOut(c) })
	private.PUT("/user/:id/archive", func(c *gin.Context) { user.Archive(c) })
	private.DELETE("/user/:id", func(c *gin.Context) { user.Delete(c) })
}
