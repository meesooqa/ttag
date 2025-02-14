package fs

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestFindFiles_File(t *testing.T) {
	logger := zaptest.NewLogger(t)
	finder := NewFinder(logger)

	tmpFile, err := os.CreateTemp("", "testfile_*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	filesChan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go finder.FindFiles(tmpFile.Name(), filesChan, &wg)

	var files []string
	for file := range filesChan {
		files = append(files, file)
	}
	wg.Wait()

	assert.Equal(t, []string{tmpFile.Name()}, files)
}

func TestFindFiles_Directory(t *testing.T) {
	logger := zaptest.NewLogger(t)
	finder := NewFinder(logger)

	tempDir, err := os.MkdirTemp("", "testdir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	file1Path := filepath.Join(tempDir, "file1.txt")
	file2Path := filepath.Join(tempDir, "file2.txt")
	err = os.WriteFile(file1Path, []byte("content1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file2Path, []byte("content2"), 0644)
	require.NoError(t, err)

	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)
	subFilePath := filepath.Join(subDir, "subfile.txt")
	err = os.WriteFile(subFilePath, []byte("subcontent"), 0644)
	require.NoError(t, err)

	filesChan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go finder.FindFiles(tempDir, filesChan, &wg)

	var files []string
	for file := range filesChan {
		files = append(files, file)
	}
	wg.Wait()

	expectedFiles := []string{file1Path, file2Path, subFilePath}
	assert.ElementsMatch(t, expectedFiles, files)
}

func TestFindFiles_NonExistent(t *testing.T) {
	logger := zaptest.NewLogger(t)
	finder := NewFinder(logger)

	nonExistentPath := "/nonexistentpath_123456789"

	filesChan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go finder.FindFiles(nonExistentPath, filesChan, &wg)

	var files []string
	for file := range filesChan {
		files = append(files, file)
	}
	wg.Wait()

	assert.Empty(t, files)
}
