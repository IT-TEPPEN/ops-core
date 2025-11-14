package dto

// FileNode represents a file or directory within a repository
type FileNode struct {
	Path string
	Type string // "file" or "dir"
}

// SelectFilesRequest represents the use case request for selecting files
type SelectFilesRequest struct {
	FilePaths []string // List of file paths to select
}
