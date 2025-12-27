package dto

// FileNode represents a file or directory within a repository
type FileNode struct {
	Path string
	Type string // "file" or "dir"
}
