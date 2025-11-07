package dto

// FileNode represents a file or directory within a repository
type FileNode struct {
	Path string `json:"path" example:"src/main.go"`
	Type string `json:"type" example:"file"` // "file" or "dir"
}

// ListFilesResponse represents the response for listing repository files
type ListFilesResponse struct {
	Files []FileNode `json:"files"`
}

// SelectFilesRequest represents the request body for selecting files
type SelectFilesRequest struct {
	FilePaths []string `json:"filePaths" binding:"required,dive,required" example:"[\"README.md\", \"docs/adr/0001.md\"]"` // List of file paths to select
}

// SelectFilesResponse represents the success response for selecting files
type SelectFilesResponse struct {
	Message       string `json:"message" example:"Files selected successfully"`
	RepoID        string `json:"repoId" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	SelectedFiles int    `json:"selectedFiles" example:"2"`
}

// GetMarkdownResponse represents the response containing selected Markdown content
type GetMarkdownResponse struct {
	RepoID  string `json:"repoId" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Content string `json:"content" example:"# Project Title\\n\\n## ADR 1\\n..."` // Concatenated Markdown content
}
