package models

// CopyPath copies a file or a folder
type CopyPath struct {
	SourceName string `json:"sourceName" binding:"required"`
	SourcePath string `json:"sourcePath" binding:"required"`
	TargetPath string `json:"targetPath" binding:"required"`
}
