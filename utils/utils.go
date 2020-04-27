package utils

import (
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func MinutesToMilliseconds(min int) int {
	return min * 60000
}

func HoursToMilliseconds(hours int) int {
	return hours * 3600000
}
