package value_object

import "errors"

// DocumentSource represents the source location of a document in a Git repository.
// It combines the file path and commit hash to uniquely identify a document version's source.
type DocumentSource struct {
	filePath   FilePath
	commitHash CommitHash
}

// NewDocumentSource creates a new DocumentSource from a file path and commit hash.
func NewDocumentSource(filePath FilePath, commitHash CommitHash) (DocumentSource, error) {
	if filePath.IsEmpty() {
		return DocumentSource{}, errors.New("file path cannot be empty")
	}
	if commitHash.IsEmpty() {
		return DocumentSource{}, errors.New("commit hash cannot be empty")
	}

	return DocumentSource{
		filePath:   filePath,
		commitHash: commitHash,
	}, nil
}

// FilePath returns the file path.
func (d DocumentSource) FilePath() FilePath {
	return d.filePath
}

// CommitHash returns the commit hash.
func (d DocumentSource) CommitHash() CommitHash {
	return d.commitHash
}

// Equals checks if two DocumentSources are equal.
func (d DocumentSource) Equals(other DocumentSource) bool {
	return d.filePath.Equals(other.filePath) && d.commitHash.Equals(other.commitHash)
}

// String returns a string representation of the document source.
func (d DocumentSource) String() string {
	return d.filePath.String() + "@" + d.commitHash.Short()
}
