package desktop

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"rakoon/rakoon-back/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetDirectory returns a directory's content
func GetDirectory(c *gin.Context) {
	const basePath = "/Users/thomasraout/Downloads"
	var path string = basePath + c.Query("path")
	var fileInfos []os.FileInfo
	var directory []models.FileDescriptor
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

	var fileDescriptor models.FileDescriptor
	for _, fileInfo := range fileInfos {
		var name = fileInfo.Name()

		if name[0] == '.' {
			continue
		}

		fileDescriptor.Name = name
		var trimmedName = trimName(name)
		fileDescriptor.TrimmedName = trimmedName
		if fileInfo.IsDir() {
			fileDescriptor.Type = "directory"
		} else {
			var extension = strings.ToLower(filepath.Ext(name))
			if extension == ".png" || extension == ".jpg" || extension == ".svg" {
				fileDescriptor.Type = "image"
			} else if extension == "mp4" {
				fileDescriptor.Type = "video"
			} else {
				fileDescriptor.Type = "file"
			}
		}
		directory = append(directory, fileDescriptor)
	}
	// directory.Directories = directories
	// directory.Files = files
	c.JSON(200, directory)
	return
}

func trimName(name string) string {
	if len(name) > 13 {
		return name[0:12] + "..."
	} else {
		return name
	}
}
