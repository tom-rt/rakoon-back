package navigation

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// GetDirectory returns a directory's content
func GetDirectory(c *gin.Context) {
	var path string = c.Query("path")
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
	fmt.Println("##### DIRECTORIES ######")
	fmt.Println(directories)
	fmt.Println("##### FILES ######")
	fmt.Println(files)
	fmt.Println("###########")
	return
}
