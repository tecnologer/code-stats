package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultDirPermissions  = 0o0755 // rwxr-xr-x
	DefaultFilePermissions = 0o0644 // rw-r--r--
	StatsDirectoryPath     = ".stats"
)

func ListFilesFromDir(dir string) []string {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", path, err)
		}

		if info.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func ReadContent(filePath string) ([]byte, error) {
	// check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", filePath)
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	return fileContent, nil
}

func Write(content []byte) error {
	err := CreateStatsFolderIfNotExists(StatsDirectoryPath)
	if err != nil {
		return fmt.Errorf("failed to create stats folder: %w", err)
	}

	fileName := fmt.Sprintf("%s/%s.json", StatsDirectoryPath, time.Now().UTC().Format("200601021_150405"))

	err = os.WriteFile(fileName, content, DefaultFilePermissions)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", fileName, err)
	}

	return nil
}

func CreateStatsFolderIfNotExists(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, DefaultDirPermissions)
		if err != nil {
			return fmt.Errorf("failed to create stats folder: %w", err)
		}
	}

	return nil
}
