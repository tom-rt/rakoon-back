package models

// PathRename renames a file or a folder
type PathRename struct {
	Name         string `json:"name" binding:"required"`
	NewPath      string `json:"newPath" binding:"required"`
	OriginalPath string `json:"originalPath" binding:"required"`
}
