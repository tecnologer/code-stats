package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
