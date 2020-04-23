package routes

import "github.com/gin-gonic/gin"
import "rakoon/user-service/authentication"

// InitRoutes calls the routes init
func InitRoutes(r *gin.Engine) {
	initUtilsRoutes(r)
	initAuthRoutes(r)
}

func initUtilsRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func initAuthRoutes(r *gin.Engine) {

	r.POST("/connect", func(c *gin.Context) {
		authentication.Connect(c)
	})

	r.POST("/subscribe", func(c *gin.Context) {
		authentication.Subscribe(c)		
	})

}
