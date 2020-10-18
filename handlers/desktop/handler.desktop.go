package desktop

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"rakoon/rakoon-back/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// DeletePath deletes at a given path
func DeletePath(c *gin.Context) {
	var rootPath = os.Getenv("ROOT_PATH")
	var pathDelete models.PathDelete

	err := c.BindJSON(&pathDelete)

	// Check formatting
	if err != nil {
		c.JSON(400, gin.H{"Incorrect input data": err.Error()})
		return
	}

	var path string = rootPath + pathDelete.Path

	err = os.RemoveAll(path)
	if err != nil {
		c.JSON(500, gin.H{"Could not remove": err.Error()})
		return
	}

	c.JSON(200, pathDelete.Path)
	return
}

// RenamePath renames a file or a directory
func RenamePath(c *gin.Context) {
	var rootPath = os.Getenv("ROOT_PATH")
	var fileRename models.PathRename
	err := c.BindJSON(&fileRename)

	// Check formatting
	if err != nil {
		c.JSON(400, gin.H{"Incorrect input data": err.Error()})
		return
	}

	var name string = fileRename.Name
	var originalPath string = rootPath + fileRename.OriginalPath
	var newPath string = rootPath + fileRename.NewPath

	err = os.Rename(originalPath, newPath)
	if err != nil {
		c.JSON(500, gin.H{"Could not rename file": err.Error()})
		return
	}

	c.JSON(201, name)
	return
}

// CreateFolder returns a directory's content
func CreateFolder(c *gin.Context) {
	var rootPath = os.Getenv("ROOT_PATH")
	var folder models.Folder
	err := c.BindJSON(&folder)

	// Check formatting
	if err != nil {
		c.JSON(400, gin.H{"Incorrect input data": err.Error()})
		return
	}

	var path string = rootPath + folder.Path
	var folderPath string = path[strings.LastIndex(folder.Name, "/")+1:]

	err = os.Mkdir(folderPath, 0755)
	if err != nil {
		c.JSON(500, gin.H{"Could not create file": err.Error()})
		return
	}

	c.JSON(201, folder.Name)
	return
}

// ServeFile returns a directory's content
func ServeFile(c *gin.Context) {
	var rootPath = os.Getenv("ROOT_PATH")
	var path string = rootPath + c.Query("path")
	var fileName string = path[strings.LastIndex(path, "/")+1:]

	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	m := http.DetectContentType(b[:512])

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, m, b)
}

// GetDirectory returns a directory's content
func GetDirectory(c *gin.Context) {
	var rootPath = os.Getenv("ROOT_PATH")
	var path string = rootPath + c.Query("path")
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
	c.JSON(200, directory)
	return
}

func trimName(name string) string {
	if len(name) > 15 {
		return name[0:13] + "..."
	}
	return name
}
