package main

import (
    "bufio"
    "flag"
    "fmt"
    "gopkg.in/yaml.v3"
    "io"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

// Config represents the structure of the YAML configuration file
type Config struct {
    Output         string   `yaml:"output"`
    InputDirs      []string `yaml:"input_dirs"`
    ExcludeFolders []string `yaml:"exclude_folders"`
    ExcludeFiles   []string `yaml:"exclude_files"`
    // Optional Enhancements:
    // IncludeFileTypes []string `yaml:"include_file_types"`
    // ExcludeFileTypes []string `yaml:"exclude_file_types"`
    // LogLevel        string   `yaml:"log_level"`
}

// ./onefile -config config.yaml

func main() {
    // Define and parse command-line flags
    configPath := flag.String("config", "config.yaml", "Path to the YAML configuration file")
    flag.Parse()

    // Read and parse the YAML configuration
    config, err := parseConfig(*configPath)
    if err != nil {
        fmt.Printf("Error parsing configuration: %v\n", err)
        os.Exit(1)
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

    // Create or truncate the output file
    outFile, err := os.Create(config.Output)
    if err != nil {
        fmt.Printf("Error creating output file '%s': %v\n", config.Output, err)
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

                // Optional: File type filtering
                // Uncomment and modify if needed
                /*
                    include := true
                    if len(config.IncludeFileTypes) > 0 {
                        include = false
                        for _, ext := range config.IncludeFileTypes {
                            if filepath.Ext(path) == ext {
                                include = true
                                break
                            }
                        }
                    }
                    if !include {
                        return nil
                    }

                    for _, ext := range config.ExcludeFileTypes {
                        if filepath.Ext(path) == ext {
                            return nil
                        }
                    }
                */

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

    fmt.Printf("All files have been combined into '%s'\n", config.Output)
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
