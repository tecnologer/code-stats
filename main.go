package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tecnologer.net/code-stats/chart"
	"tecnologer.net/code-stats/command"
	"tecnologer.net/code-stats/models"
)

const statsFolder = ".stats"

var (
	inputFilePaths   = flag.String("input", ".stats", "Path to the input file, separated by commas, could be used with stdin")
	imitDir          = flag.String("omit-dir", ".idea,vendor,.stats", "Directories to omit from the stats")
	onlyCompareInput = flag.Bool("only-compare-input", false, "Only compare the input files, do not calculate the current stats")
	drawChart        = flag.Bool("draw-chart", false, "Draw chart")
	languages        = flag.String("languages", "go", "Languages to include in the chart, require at least one and --draw-chart")
	showVersion      = flag.Bool("version", false, "Show version")
	version          string
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println(version)

		os.Exit(0)
	}

	stats, err := getJSONContent()
	if err != nil {
		log.Fatal(err)
	}

	if !*onlyCompareInput {
		currentStats, err := calculateCurrentStats()
		if err != nil {
			log.Fatal(err)
		}

		if stats != nil {
			stats.Merge(currentStats)
		} else {
			stats = currentStats
		}
	}

	fmt.Println("Stats collected successfully")

	if *drawChart {
		langs := strings.Split(*languages, ",")

		err = chart.Draw(stats, langs...)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Chart generated successfully")
	}
}

func calculateCurrentStats() (*models.StatsCollection, error) {
	output, err := command.RunSCC(*imitDir)
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

	err = writeFile(output)
	if err != nil {
		log.Fatal(err)
	}

	collection := models.NewCollection()
	collection.Add(key, statsData)

	return collection, nil
}

func writeFile(content []byte) error {
	err := createStatsFolderIfNotExists()
	if err != nil {
		return fmt.Errorf("failed to create stats folder: %w", err)
	}

	fileName := fmt.Sprintf("%s/%s.json", statsFolder, time.Now().UTC().Format(time.DateOnly))

	err = os.WriteFile(fileName, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", fileName, err)
	}

	return nil

}

func createStatsFolderIfNotExists() error {
	_, err := os.Stat(statsFolder)
	if os.IsNotExist(err) {
		err = os.Mkdir(statsFolder, 0755)
		if err != nil {
			return fmt.Errorf("failed to create stats folder: %w", err)
		}
	}

	return nil
}

func getJSONContent() (*models.StatsCollection, error) {
	dataFromFiles, err := readJSONContentFromFiles()
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

func readJSONContentFromFiles() (*models.StatsCollection, error) {
	if *inputFilePaths == "" {
		return nil, nil
	}

	paths := strings.Split(*inputFilePaths, ",")
	filePaths := make([]string, 0, len(paths))

	stats := models.NewCollection()

	for _, filePath := range paths {
		stat, _ := os.Stat(filePath)
		if stat.IsDir() {
			dirFilePaths := listFilesFromDir(filePath)
			filePaths = append(filePaths, dirFilePaths...)

			continue
		}

		filePaths = append(filePaths, filePath)
	}

	for _, filePath := range filePaths {
		jsonContent, err := readJSONContentFromFile(filePath)
		if err != nil {
			return nil, err
		}

		var statsByFile *models.StatsCollection

		err = json.Unmarshal(jsonContent, &statsByFile)
		if err != nil {
			return nil, err
		}

		stats.Merge(statsByFile)
	}

	return stats, nil
}

func readJSONContentFromFile(filePath string) ([]byte, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	return fileContent, nil
}

func readJSONContentFromStdin() (*models.StatsCollection, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, nil
	}

	scanner := bufio.NewScanner(os.Stdin)

	var text bytes.Buffer

	for scanner.Scan() {
		text.Write(scanner.Bytes())
		text.WriteString("\n")
	}

	if text.Len() == 0 {
		return nil, nil
	}

	var stats *models.StatsCollection

	err := json.Unmarshal(text.Bytes(), &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON content from stdin: %w", err)
	}

	return stats, nil
}

func listFilesFromDir(dir string) []string {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
