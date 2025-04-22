package model

// FileNode represents a file or directory within a repository.
type FileNode struct {
	path  string // Path relative to the repository root
	type_ string // Type, e.g., "file" or "dir"
}

// NewFileNode creates a new FileNode with the given path and type.
func NewFileNode(path, type_ string) *FileNode {
	return &FileNode{
		path:  path,
		type_: type_,
	}
}

// ReconstructFileNode reconstructs a FileNode from persistence or external data.
func ReconstructFileNode(path, type_ string) *FileNode {
	return &FileNode{
		path:  path,
		type_: type_,
	}
}

// Path returns the file path relative to the repository root.
func (f *FileNode) Path() string {
	return f.path
}

// Type returns the type of the file node (e.g., "file" or "dir").
func (f *FileNode) Type() string {
	return f.type_
}
