package file //nolint:testpackage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFilesFromDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	// Create sample files and directories
	_, err := os.Create(filepath.Join(dir, "file1.txt"))
	require.NoError(t, err)

	_, err = os.Create(filepath.Join(dir, "file2.txt"))
	require.NoError(t, err)

	err = os.Mkdir(filepath.Join(dir, "subdir"), DefaultFilePermissions)
	require.NoError(t, err)

	expected := []string{
		filepath.Join(dir, "file1.txt"),
		filepath.Join(dir, "file2.txt"),
	}

	files := ListFilesFromDir(dir)
	assert.Equal(t, expected, files)
}

func TestReadContent(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "testfile.txt")
	content := []byte("hello world")

	err := os.WriteFile(filePath, content, DefaultDirPermissions)
	require.NoError(t, err)

	readContent, err := ReadContent(filePath)
	require.NoError(t, err)

	assert.Equal(t, content, readContent)
}

func TestWrite(t *testing.T) {
	t.Parallel()

	content := []byte("test content")

	require.NoError(t, Write(content))

	// Verify file is written
	files, err := os.ReadDir(StatsDirectoryPath)
	require.NoError(t, err)
	assert.NotEmpty(t, files)

	// Cleanup
	err = os.RemoveAll(StatsDirectoryPath)
	require.NoError(t, err)
}

func TestCreateStatsFolderIfNotExists(t *testing.T) {
	t.Parallel()

	dir := filepath.Join(os.TempDir(), "stats_test")

	// Ensure directory does not exist
	err := os.RemoveAll(dir)
	require.NoError(t, err)

	err = CreateStatsFolderIfNotExists(dir)
	require.NoError(t, err)

	// Check if directory exists
	_, err = os.Stat(dir)
	require.NoError(t, err)

	// Cleanup
	err = os.RemoveAll(dir)
	require.NoError(t, err)
}
