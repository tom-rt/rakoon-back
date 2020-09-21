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

type File struct {
	Name string `uri:"name" binding:"required"`
}

// ServeFile returns a directory's content
func ServeFile(c *gin.Context) {
	fmt.Println("SERVE")
	// c.File("/home/thomas/programming/perso/rakoon-back/main.go")
	const path = "/home/thomas/programming/perso/rakoon-back/main.go"

	var f File
	err := c.ShouldBindUri(&f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {
		fmt.Println(err)
	}
	fmt.Println("ICI")

	// %%%%%%%%%%
	// dst := fmt.Sprintf("%s/%s", base, f.Name)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	m := http.DetectContentType(b[:512])

	// %%%%%%%%%%

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+f.Name)
	fmt.Println("ICI", m, b)
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
	if len(name) > 13 {
		return name[0:12] + "..."
	} else {
		return name
	}
}
