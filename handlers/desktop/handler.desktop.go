package desktop

import (
	"io/ioutil"
	"log"
	"os"
	"rakoon/rakoon-back/models"

	"github.com/gin-gonic/gin"
)

// GetDirectory returns a directory's content
func GetDirectory(c *gin.Context) {
	const basePath = "/home/thomas"
	var path string = basePath + c.Query("path")
	var fileInfos []os.FileInfo
	var directories []string
	var files []string
	if len(path) <= 0 {
		c.JSON(401, gin.H{
			"message": "No path specified.",
		})
		return
	}
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			directories = append(directories, fileInfo.Name())
		} else {
			files = append(files, fileInfo.Name())
		}
	}
	var directory = new(models.Directory)
	directory.Directories = directories
	directory.Files = files
	c.JSON(200, directory)
	return
}
