package schema

import "time"

// FileNode represents a file or directory within a repository in API responses
type FileNode struct {
	Path       string  `json:"path" example:"src/main.go"`
	Type       string  `json:"type" example:"file"` // "file" or "dir"
	SHA        *string `json:"sha,omitempty" example:"abc1234567890def"`
	Size       *int    `json:"size,omitempty" example:"1024"` // File size in bytes (only for files)
	CommitHash *string `json:"commit_hash,omitempty" example:"def4567890123abc"`
}

// ListFilesResponse represents the API response for listing repository files
type ListFilesResponse struct {
	Files          []FileNode `json:"files"`
	LatestCommit   *string    `json:"latest_commit,omitempty" example:"abc1234567890def"` // Latest commit hash of the repository
	LatestCommitAt *time.Time `json:"latest_commit_at,omitempty" example:"2025-04-22T10:00:00Z"`
}

// GetFileContentsResponse represents the API response containing file content
type GetFileContentsResponse struct {
	RepoID     string  `json:"repoId" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	FilePath   string  `json:"filePath" example:"README.md"`
	Content    string  `json:"content" example:"# Project Title\n\nThis is the README content..."`
	CommitHash *string `json:"commit_hash,omitempty" example:"abc1234567890def"` // Commit hash of the file
	SHA        *string `json:"sha,omitempty" example:"def4567890123abc"`         // SHA of the file content
}

// FileCommitInfo represents commit information for a file
type FileCommitInfo struct {
	CommitHash string    `json:"commit_hash" example:"abc1234567890def"`
	Message    string    `json:"message" example:"Updated documentation"`
	Author     string    `json:"author" example:"John Doe"`
	AuthorEmail string   `json:"author_email" example:"john@example.com"`
	Date       time.Time `json:"date" example:"2025-04-22T10:00:00Z"`
}

// GetFileHistoryResponse represents the API response for file commit history
type GetFileHistoryResponse struct {
	FilePath string           `json:"file_path" example:"docs/procedure.md"`
	Commits  []FileCommitInfo `json:"commits"`
}
