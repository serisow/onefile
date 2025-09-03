package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of the YAML configuration file
type Config struct {
	Output         string   `yaml:"output"`
	OutputDir      string   `yaml:"output_dir"`
	InputDirs      []string `yaml:"input_dirs"`
	ExcludeFolders []string `yaml:"exclude_folders"`
	ExcludeFiles   []string `yaml:"exclude_files"`
}

// ./onefile -config config.yaml
func main() {
	// Define and parse command-line flags
	configPath := flag.String("config", "config.yaml", "Path to the YAML configuration file")
	// Add this line with other flag definitions
	split := flag.Bool("split", false, "Split the combined output file into two nearly equal halves")

	flag.Parse()

	// Read and parse the YAML configuration
	config, err := parseConfig(*configPath)
	if err != nil {
		fmt.Printf("Error parsing configuration: %v\n", err)
		os.Exit(1)
	}

	// Get the directory of the configuration file to resolve relative paths
	configDir := filepath.Dir(*configPath)

	// Resolve relative paths in the configuration
	for i, dir := range config.InputDirs {
		if !filepath.IsAbs(dir) {
			config.InputDirs[i] = filepath.Join(configDir, dir)
		}
	}
	for i, folder := range config.ExcludeFolders {
		if !filepath.IsAbs(folder) {
			config.ExcludeFolders[i] = filepath.Join(configDir, folder)
		}
	}
	for i, file := range config.ExcludeFiles {
		if !filepath.IsAbs(file) {
			config.ExcludeFiles[i] = filepath.Join(configDir, file)
		}
	}

	// Validate input directories
	if len(config.InputDirs) == 0 {
		fmt.Println("Error: No input directories specified in the configuration.")
		os.Exit(1)
	}

	// Set default output file if not specified
	if config.Output == "" {
		config.Output = "combined.txt"
	}

	// Construct the full output path
	var outputPath string
	if config.OutputDir != "" {
		// If the output directory is not an absolute path, join it with the config file's directory
		if !filepath.IsAbs(config.OutputDir) {
			config.OutputDir = filepath.Join(configDir, config.OutputDir)
		}
		// Create the output directory if it doesn't exist
		if err := os.MkdirAll(config.OutputDir, os.ModePerm); err != nil {
			fmt.Printf("Error creating output directory '%s': %v\n", config.OutputDir, err)
			os.Exit(1)
		}
		outputPath = filepath.Join(config.OutputDir, config.Output)
	} else {
		// If no output directory is specified, place the output file in the same directory as the config file
		outputPath = filepath.Join(filepath.Dir(*configPath), config.Output)
	}

	// Create or truncate the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file '%s': %v\n", outputPath, err)
		os.Exit(1)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// Mutex to synchronize writes to the writer
	var writerMutex sync.Mutex

	// Use a WaitGroup to handle concurrency
	var wg sync.WaitGroup
	fileChan := make(chan string, 100) // Buffered channel to hold file paths

	// Start a fixed number of worker goroutines
	numWorkers := 4 // Adjust based on your CPU cores
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for filePath := range fileChan {
				fmt.Printf("Worker %d processing file: %s\n", workerID, filePath)
				err := appendFileContent(filePath, writer, &writerMutex)
				if err != nil {
					fmt.Printf("Error reading file '%s': %v\n", filePath, err)
				} else {
					// Optionally, add a separator between files
					separator := fmt.Sprintf("\n--- End of %s ---\n\n", filePath)
					writerMutex.Lock()
					_, err = writer.WriteString(separator)
					if err != nil {
						fmt.Printf("Error writing separator for '%s': %v\n", filePath, err)
					}
					writerMutex.Unlock()
				}
			}
		}(i + 1)
	}

	// Iterate over input directories
	for _, dir := range config.InputDirs {
		// Check if the directory exists
		fi, err := os.Stat(dir)
		if err != nil {
			fmt.Printf("Error accessing directory '%s': %v\n", dir, err)
			continue
		}
		if !fi.IsDir() {
			fmt.Printf("Skipping '%s': not a directory.\n", dir)
			continue
		}

		// Walk through the directory
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Error accessing path '%s': %v\n", path, err)
				return nil // Skip this file/directory but continue walking
			}

			// Check for excluded folders
			for _, exclDir := range config.ExcludeFolders {
				if isSubPath(exclDir, path) && info.IsDir() {
					fmt.Printf("Excluding directory: %s\n", path)
					return filepath.SkipDir
				}
			}

			if !info.IsDir() {
				// Check for excluded files
				exclude := false
				for _, exclFile := range config.ExcludeFiles {
					if filepath.Clean(exclFile) == filepath.Clean(path) {
						fmt.Printf("Excluding file: %s\n", path)
						exclude = true
						break
					}
				}
				if exclude {
					return nil
				}
				// Send the file path to the channel for processing
				fileChan <- path
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error walking through directory '%s': %v\n", dir, err)
		}
	}

	// Close the channel and wait for all workers to finish
	close(fileChan)
	wg.Wait()

	fmt.Printf("All files have been combined into '%s'\n", outputPath)

	// Handle splitting if the -split flag is set
	if *split {
		err := splitFile(outputPath)
		if err != nil {
			fmt.Printf("Error splitting the file '%s': %v\n", outputPath, err)
			os.Exit(1)
		}
		ext := filepath.Ext(outputPath)
		baseName := strings.TrimSuffix(outputPath, ext)
		fmt.Printf("The file '%s' has been split into '%s_part1%s' and '%s_part2%s'\n",
			outputPath,
			baseName,
			ext,
			baseName,
			ext)
	}
}

// parseConfig reads and parses the YAML configuration file
func parseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// isSubPath checks if sub is a subdirectory or the same as base
func isSubPath(base, sub string) bool {
	rel, err := filepath.Rel(base, sub)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, "..") && rel != "")
}

// appendFileContent reads the content of the given file and writes it to the writer
func appendFileContent(filePath string, writer *bufio.Writer, mutex *sync.Mutex) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Optionally, write the file name as a header
	header := fmt.Sprintf("\n--- %s ---\n", filePath)
	mutex.Lock()
	_, err = writer.WriteString(header)
	if err != nil {
		mutex.Unlock()
		return err
	}
	mutex.Unlock()

	// Copy file content
	// To ensure thread safety, lock the writer during the copy
	mutex.Lock()
	_, err = io.Copy(writer, file)
	mutex.Unlock()
	return err
}

// splitFile splits the given file into two nearly equal halves based on the number of lines.
func splitFile(filePath string) error {
	// Open the combined file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open the combined file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading the combined file: %v", err)
	}

	totalLines := len(lines)
	if totalLines == 0 {
		return fmt.Errorf("the combined file is empty")
	}

	// Determine the split point
	splitPoint := totalLines / 2

	// Define the names for the split files
	ext := filepath.Ext(filePath)
	baseName := strings.TrimSuffix(filePath, ext)
	part1 := fmt.Sprintf("%s_part1%s", baseName, ext)
	part2 := fmt.Sprintf("%s_part2%s", baseName, ext)

	// Write the first half to part1
	err = writeLinesToFile(part1, lines[:splitPoint])
	if err != nil {
		return fmt.Errorf("failed to write to %s: %v", part1, err)
	}

	// Write the second half to part2
	err = writeLinesToFile(part2, lines[splitPoint:])
	if err != nil {
		return fmt.Errorf("failed to write to %s: %v", part2, err)
	}

	return nil
}

// writeLinesToFile writes the given lines to the specified file.
func writeLinesToFile(filePath string, lines []string) error {
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file %s: %v", filePath, err)
		}
	}

	// Flush the buffer to ensure all data is written
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer for file %s: %v", filePath, err)
	}

	return nil
}
