package torrent

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// Download uploads a file
func Download(c *gin.Context) {
	var rootPath = os.Getenv("ROOT_PATH")
	var path string = rootPath + c.PostFormArray("path")[0]

	file, err := c.FormFile("file")
	src, err := file.Open()
	if err != nil {
		c.JSON(401, err.Error())
		return
	}
	defer src.Close()

	target := fmt.Sprintf("%s/%s", path, file.Filename)
	out, err := os.Create(target)
	if err != nil {
		c.JSON(402, err.Error())
		return
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	c.JSON(201, "File(s) uploaded.")
	return
}
