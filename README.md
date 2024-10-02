## 📚 README.md

```markdown
# 🗂️ OneFile - Combine Multiple Files into a Single Text File

OneFile is a powerful and convenient command-line tool written in Go that allows you to recursively combine multiple files from specified directories into a single text file. It's perfect for aggregating code, logs, or any text-based files, making it easy to copy and paste the combined content into platforms like ChatGPT or any other text-based application.

## 🚀 Features

- **Recursive File Traversal:** Automatically walks through provided directories and their subdirectories to locate all files.
- **YAML Configuration:** Easily configure input directories, exclusions, and output file names using a YAML file.
- **Exclusion Lists:** Exclude specific folders and files from the aggregation process.
- **Concurrency:** Utilizes multiple goroutines to process files concurrently for faster execution.
- **Thread-Safe Writing:** Ensures data integrity with synchronized write operations.
- **Customizable Output:** Specify the name and location of the combined output file.
- **Easy to Use:** Simple command-line interface with clear configuration options.

## 🛠️ Installation

### Prerequisites

- **Go:** Ensure you have Go installed on your system. You can download it from [golang.org](https://golang.org/dl/).

### Clone the Repository

```bash
git clone https://github.com/serisow/onefile.git
cd onefile
```

### Build the Executable

```bash
go build -o onefile main.go
```

This command compiles the Go source code and generates an executable named `onefile` (or `onefile.exe` on Windows) in your current directory.

## 📋 Usage

OneFile operates based on a YAML configuration file. This approach provides flexibility and ease of use.

### 1. Create a Configuration File

Create a `config.yaml` file in your project directory with the following structure:

```yaml
output: combined.txt
input_dirs:
  - /path/to/dir1
  - /path/to/dir2
exclude_folders:
  - /path/to/dir1/exclude_this_folder
  - /path/to/dir2/another_excluded_folder
exclude_files:
  - /path/to/dir1/file_to_exclude.txt
  - /path/to/dir2/another_file_to_exclude.md
# Optional Enhancements:
# include_file_types:
#   - .txt
#   - .md
# exclude_file_types:
#   - .log
# log_level: info
```

#### Configuration Fields

- **output**: *(string)* Name of the output file. Defaults to `combined.txt` if not specified.
- **input_dirs**: *(list of strings)* Paths to the input directories containing files to be combined.
- **exclude_folders**: *(list of strings)* Paths to folders that should be excluded from processing.
- **exclude_files**: *(list of strings)* Specific file paths to exclude from the aggregation.
- **include_file_types**: *(list of strings)* *(Optional)* Specify file extensions to include.
- **exclude_file_types**: *(list of strings)* *(Optional)* Specify file extensions to exclude.
- **log_level**: *(string)* *(Optional)* Define the logging verbosity (e.g., `info`, `warning`, `error`).

### 2. Run the Tool

Execute the `onefile` binary with the configuration file:

```bash
./onefile -config config.yaml
```

**Windows:**

```cmd
onefile.exe -config config.yaml
```

### 3. Output

After running, the tool will generate the `combined.txt` file (or your specified output file) containing the concatenated contents of all processed files, separated by headers and footers for clarity.

## 🔧 Configuration Example

Here's a sample `config.yaml` for reference:

```yaml
output: all_code_combined.txt
input_dirs:
  - /home/user/projects/project1
  - /home/user/projects/project2
exclude_folders:
  - /home/user/projects/project1/vendor
  - /home/user/projects/project2/tmp
exclude_files:
  - /home/user/projects/project1/main_test.go
  - /home/user/projects/project2/README.md
include_file_types:
  - .go
  - .py
exclude_file_types:
  - .log
log_level: info
```

## 🧩 Additional Features

While OneFile is already robust, here are some additional features you might consider implementing:

- **Progress Indicator:** Integrate a progress bar to provide visual feedback during processing.
- **Logging Enhancements:** Use a logging library to support different logging levels and outputs.
- **Dry-Run Mode:** Allow users to simulate the aggregation process without writing to the output file.
- **Compression Option:** Offer the ability to compress the combined output file.
- **GUI Wrapper:** Develop a simple graphical interface for users who prefer not to use the command line.

## 📝 Contributing

Contributions are welcome! If you'd like to contribute to OneFile, please follow these steps:

1. **Fork the Repository:** Click the "Fork" button at the top right of the repository page.
2. **Clone Your Fork:**

   ```bash
   git clone https://github.com/serisow/onefile.git
   cd onefile
   ```

3. **Create a Feature Branch:**

   ```bash
   git checkout -b feature/YourFeature
   ```

4. **Commit Your Changes:**

   ```bash
   git commit -m "Add Your Feature"
   ```

5. **Push to Your Fork:**

   ```bash
   git push origin feature/YourFeature
   ```

6. **Create a Pull Request:** Go to your fork on GitHub and click the "Compare & pull request" button.

## 🐛 Troubleshooting

### Common Issues

- **Short Write Error:**

  **Error Message:**
  ```
  Error reading file '/path/to/file': short write
  panic: runtime error: slice bounds out of range [X:Y]
  ```

  **Cause:**
  Concurrent writes to the output file without proper synchronization.

  **Solution:**
  Ensure you are using the latest version of OneFile with proper mutex synchronization. If the issue persists, please open an issue on GitHub with detailed logs.

### Running with Race Detector

To check for race conditions, run the tool with Go's race detector:

```bash
go run -race main.go -config config.yaml
```

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 📝 Acknowledgements

- Built with [Go](https://golang.org/)
- YAML configuration using [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3)

## 📫 Contact

For any questions or feedback, please open an issue on the [GitHub repository](https://github.com/serisow/onefile) or contact [drupaliste@gmail.com](mailto:drupaliste@gmail.com).

```

---

## 🗂️ .gitignore

Create a `.gitignore` file in the root of your project directory to exclude unnecessary files from your Git repository. Here's a recommended `.gitignore` tailored for Go projects, along with some additional exclusions specific to your tool:

```gitignore
# Binaries
/bin/
/*.exe
/*.exe~
/*.dll
/*.so
/*.dylib

# Output files
/combined.txt
/*.log



# Vendor directory (if using Go modules)
/vendor/

# Go workspace file
/go.work

# IDE and editor directories
/.idea/
/.vscode/
/*.swp
.DS_Store

# Configuration files (if any sensitive data)
config.yaml

```

### Explanation of `.gitignore` Entries

- **Binaries:**
  - Exclude compiled binaries and dynamic libraries to avoid cluttering the repository.

- **Output Files:**
  - Exclude `combined.txt` and any `.log` files generated during execution.

- **Test and Coverage Files:**
  - Exclude test binaries and coverage reports.

- **Vendor Directory:**
  - If you're using Go modules with a `vendor` directory, it's often excluded unless necessary.

- **IDE and Editor Directories:**
  - Exclude configuration directories for popular IDEs like IntelliJ IDEA (`.idea/`) and Visual Studio Code (`.vscode/`), as well as swap files.

- **OS-specific Files:**
  - Exclude system-generated files like `.DS_Store` for macOS and `Thumbs.db` for Windows.

- **Configuration Files:**
  - Exclude `config.yaml` if it contains sensitive information. If your configuration doesn't contain sensitive data, you might choose to include it.

- **Temporary Files:**
  - Exclude temporary files and directories to keep the repository clean.

