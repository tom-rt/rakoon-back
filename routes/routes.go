package routes

import (
	"rakoon/rakoon-back/handlers/authentication"
	"rakoon/rakoon-back/handlers/desktop"
	"rakoon/rakoon-back/handlers/user"
	"rakoon/rakoon-back/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter calls the routes init
func SetupRouter() *gin.Engine {
	router := gin.New()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	router.Use(cors.New(config))

	// Public routes
	public := router.Group("/v1")

	public.POST("/user", func(c *gin.Context) { user.Create(c) })
	public.POST("/user/login", func(c *gin.Context) { user.Connect(c) })
	public.POST("/refresh/token", func(c *gin.Context) { authentication.RefreshToken(c) })

	// Private Routes
	private := router.Group("/v1")
	private.Use(middleware.JwtHandling)

	private.GET("/user/:id", func(c *gin.Context) { user.Get(c) })
	private.GET("/list/directory", func(c *gin.Context) { desktop.GetDirectory(c) })

	// private.GET("/local/file", func(c *gin.Context) { c.File("local/file.go") })
	private.GET("/filez", func(c *gin.Context) { desktop.ServeFile(c) })

	// var fs http.FileSystem = // ...
	// router.GET("/fs/file", func(c *gin.Context) {
	// 	c.FileFromFS("fs/file.go", fs)
	// })

	private.PUT("/user/:id", func(c *gin.Context) { user.Update(c) })
	private.PUT("/user/:id/password", func(c *gin.Context) { user.UpdatePassword(c) })
	private.PUT("/user/:id/logout", func(c *gin.Context) { user.LogOut(c) })
	private.PUT("/user/:id/archive", func(c *gin.Context) { user.Archive(c) })
	private.DELETE("/user/:id", func(c *gin.Context) { user.Delete(c) })

	return router
}
