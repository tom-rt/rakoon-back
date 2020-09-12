package models

// FileDescriptor type represents the content of a directory
type FileDescriptor struct {
	TrimmedName string `json:"trimmedName"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}
