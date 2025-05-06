package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileNode(t *testing.T) {
	// テスト：FileNodeの作成が正しく行われることを確認する
	t.Run("FileNodeが正しく作成される", func(t *testing.T) {
		path := "path/to/file.md"
		fileType := "file"

		fileNode := NewFileNode(path, fileType)

		assert.NotNil(t, fileNode)
		assert.Equal(t, path, fileNode.Path())
		assert.Equal(t, fileType, fileNode.Type())
	})

	// テスト：ディレクトリタイプのFileNodeが正しく作成される
	t.Run("ディレクトリタイプのFileNodeが正しく作成される", func(t *testing.T) {
		path := "path/to/directory"
		fileType := "dir"

		fileNode := NewFileNode(path, fileType)

		assert.NotNil(t, fileNode)
		assert.Equal(t, path, fileNode.Path())
		assert.Equal(t, fileType, fileNode.Type())
	})
}

func TestReconstructFileNode(t *testing.T) {
	// テスト：FileNodeの再構築が正しく行われることを確認する
	t.Run("FileNodeが正しく再構築される", func(t *testing.T) {
		path := "path/to/file.md"
		fileType := "file"

		fileNode := ReconstructFileNode(path, fileType)

		assert.NotNil(t, fileNode)
		assert.Equal(t, path, fileNode.Path())
		assert.Equal(t, fileType, fileNode.Type())
	})
}

func TestFileNodeMethods(t *testing.T) {
	// テスト：Path()メソッドが正しいパスを返すことを確認する
	t.Run("Path()メソッドが正しいパスを返す", func(t *testing.T) {
		path := "path/to/file.md"
		fileType := "file"

		fileNode := NewFileNode(path, fileType)

		assert.Equal(t, path, fileNode.Path())
	})

	// テスト：Type()メソッドが正しいタイプを返すことを確認する
	t.Run("Type()メソッドが正しいタイプを返す", func(t *testing.T) {
		path := "path/to/file.md"
		fileType := "file"

		fileNode := NewFileNode(path, fileType)

		assert.Equal(t, fileType, fileNode.Type())
	})
}
