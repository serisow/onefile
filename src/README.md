Of course. Here is a rewritten, simplified `README.md` that keeps only the essential information and clearly describes the new usage flow.

This version is more direct, easier to scan, and focuses on the practical steps a user needs to take.

---

# 🗂️ OneFile

A command-line tool to combine multiple source files into a single text file, perfect for providing context to LLMs, creating archives, or generating documentation.

## Features

-   **Recursive Processing:** Scans specified directories to find all files.
-   **YAML Configuration:** A simple `config.yaml` to define what to include and exclude.
-   **Exclusion Rules:** Ignore specific folders (like `.git`, `node_modules`) and files.
-   **Concurrent:** Processes files in parallel for improved speed.
-   **File Splitting:** Optionally splits the final output file into two parts.

## Installation

### Option 1: Global (Recommended)

Build the binary and move it to a location in your system's `PATH` so you can run it from any directory.

1.  **Build:**
    ```bash
    go build -o onefile main.go
    ```
2.  **Move (Linux/macOS):**
    ```bash
    # Make sure ~/go/bin is in your $PATH
    mv onefile ~/go/bin/
    ```
    **Move (Windows):**
    Move `onefile.exe` to a folder that is included in your `Path` environment variable.

### Option 2: Portable

Place the compiled `onefile` binary in a known location. You will run it by providing the path to the binary.

## Usage Flow

Here is the standard workflow for using `onefile` with your project.

### 1. Set up Your Project

For this example, assume your project has this structure:

```
/path/to/my-project/
├── src/
│   ├── main.go
│   └── utils.go
├── docs/
│   └── guide.md
├── .git/
│   └── ...
└── README.md
```

### 2. Create `config.yaml`

In the root of your project (`/path/to/my-project/`), create a `config.yaml` file. This file tells `onefile` what to do.

**`config.yaml`:**
```yaml
# The name of the final combined file.
output: combined_context.txt

# The directories to scan for files.
# Paths are relative to this config file's location.
input_dirs:
  - ./src
  - ./docs

# Folders to completely ignore.
exclude_folders:
  - ./.git
  - ./vendor

# Specific files to ignore.
exclude_files:
  - ./docs/guide.md
```

### 3. Run the Tool

Open your terminal and run the `onefile` command, pointing it to your configuration file.

*   **If you installed it globally:**
    ```bash
    # Run from anywhere
    onefile -config /path/to/my-project/config.yaml
    ```

*   **If using the portable binary:**
    ```bash
    # Run from the directory containing the binary
    ./onefile -config /path/to/my-project/config.yaml
    ```

*   **To also split the output file:**
    ```bash
    onefile -config /path/to/my-project/config.yaml -split
    ```

### 4. Get the Output

The tool will create `combined_context.txt` in the same directory as your `config.yaml`. The file will contain the contents of `main.go`, `utils.go`, and `README.md`, each with a clear header indicating its original path.

## Configuration Reference

-   `output` (string): Name of the generated file. Defaults to `combined.txt`.
-   `output_dir` (string, *optional*): A directory where the output file will be saved.
-   `input_dirs` (list of strings): List of directories to scan for files. Paths can be absolute or relative to the `config.yaml` location.
-   `exclude_folders` (list of strings): List of directories to exclude from scanning.
-   `exclude_files` (list of strings): List of specific files to exclude.

*The following features are planned but not yet implemented:*
```yaml
# include_file_types: [".go", ".md"] # To only include files with these extensions.
# exclude_file_types: [".log", ".tmp"] # To exclude files with these extensions.
```