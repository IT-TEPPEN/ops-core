package schema

// FileNode represents a file or directory within a repository in API responses
type FileNode struct {
	Path string `json:"path" example:"src/main.go"`
	Type string `json:"type" example:"file"` // "file" or "dir"
}

// ListFilesResponse represents the API response for listing repository files
type ListFilesResponse struct {
	Files []FileNode `json:"files"`
}

// GetFileContentsResponse represents the API response containing file content
type GetFileContentsResponse struct {
	RepoID   string `json:"repoId" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	FilePath string `json:"filePath" example:"README.md"`
	Content  string `json:"content" example:"# Project Title\n\nThis is the README content..."`
}
