package models

// Folder type represents a folder
type Folder struct {
	Name string `json:"name" binding:"required"`
	Path string `json:"path"`
}
