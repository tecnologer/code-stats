package extractor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"tecnologer.net/code-stats/pkg/file"
	"tecnologer.net/code-stats/pkg/models"
	"tecnologer.net/code-stats/pkg/scc"
)

func ExtractFromInput(paths []string) (*models.StatsCollection, error) {
	dataFromFiles, err := readJSONContentFromFiles(paths)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON content from files: %w", err)
	}

	fromStdIn, err := readJSONContentFromStdin()
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON content from stdin: %w", err)
	}

	dataFromFiles.Merge(fromStdIn)

	return dataFromFiles, nil
}

func readJSONContentFromFiles(paths []string) (*models.StatsCollection, error) {
	if len(paths) == 0 {
		return nil, nil //nolint:nilnil
	}

	filePaths := make([]string, 0, len(paths))

	stats := models.NewCollection()

	for _, filePath := range paths {
		stat, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path %s does not exist", filePath)
		}

		if stat.IsDir() {
			dirFilePaths := file.ListFilesFromDir(filePath)
			filePaths = append(filePaths, dirFilePaths...)

			continue
		}

		filePaths = append(filePaths, filePath)
	}

	for _, filePath := range filePaths {
		jsonContent, err := file.ReadContent(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read content from file %s: %w", filePath, err)
		}

		var statsByFile *models.StatsCollection

		err = json.Unmarshal(jsonContent, &statsByFile)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON content from file %s: %w", filePath, err)
		}

		stats.Merge(statsByFile)
	}

	return stats, nil
}

func readJSONContentFromStdin() (*models.StatsCollection, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, nil //nolint:nilnil
	}

	scanner := bufio.NewScanner(os.Stdin)

	var text bytes.Buffer

	for scanner.Scan() {
		text.Write(scanner.Bytes())
		text.WriteString("\n")
	}

	if text.Len() == 0 {
		return nil, nil //nolint:nilnil
	}

	var stats *models.StatsCollection

	err := json.Unmarshal(text.Bytes(), &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON content from stdin: %w", err)
	}

	return stats, nil
}

func ExtractCurrent(omitDirs []string) (*models.StatsCollection, error) {
	output, err := scc.Run(omitDirs...)
	if err != nil {
		log.Fatal(err)
	}

	var statsData []*models.Stats

	err = json.Unmarshal(output, &statsData)
	if err != nil {
		log.Fatal(err)
	}

	key := time.Now().UTC()
	stats := map[time.Time][]*models.Stats{
		key: statsData,
	}

	output, err = json.MarshalIndent(stats, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = file.Write(output)
	if err != nil {
		log.Fatal(err)
	}

	collection := models.NewCollection()
	collection.Add(key, statsData)

	return collection, nil
}