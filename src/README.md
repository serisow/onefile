# 🗂️ OneFile

A smart command-line tool to combine multiple source files into a single text file. Perfect for providing context to LLMs, creating project archives, or generating documentation.

OneFile follows a **convention over configuration** approach. It works out-of-the-box with smart defaults, requiring minimal to zero configuration for most projects.

## How It Works

-   **Zero-Config Ready:** Run it in your project root, and it automatically scans everything.
-   **Smart Defaults:** By default, it looks for a `config-onefile.yml` file, reads from the current directory (`.`), and writes to `onefile.txt`.
-   **Auto-Exclusions:** Automatically ignores common folders (`.git`, `node_modules`, `vendor`) and files (its own binary, the config file, `go.mod`, etc.) so you don't have to.

## Installation

1.  **Build the binary:**
    ```bash
    go build -o onefile main.go
    ```

2.  **Make the binary accessible:**
    -   **Option 1: Global (Recommended)**
        Move the binary to a directory in your system's `PATH`.

        *For system-wide access on Linux/macOS (requires admin rights):*
        ```bash
        sudo mv onefile /usr/local/bin/
        ```
        *For user-specific access (common for Go developers):*
        ```bash
        # Ensure '~/go/bin' is in your $PATH
        mv onefile ~/go/bin/
        ```

    -   **Option 2: Local**
        Place the `onefile` binary in your project's root directory.

## Quick Start: Zero Configuration

For most projects, no configuration is needed.

1.  **Navigate to your project root:**
    ```bash
    cd /path/to/my-project/
    ```
2.  **Run the tool:**
    ```bash
    # If installed globally
    onefile

    # If using the local binary in your project
    ./onefile
    ```

OneFile will scan the project, apply default exclusions, and create `onefile.txt` in the root.

## Customizing with `config-onefile.yml`

To override defaults, create a `config-onefile.yml` file in your project root. You only need to specify what you want to change.

#### Example

To combine only the files in the `src` directory and name the output `app_context.txt`:

**`config-onefile.yml`:**
```yaml
output: app_context.txt
input_dirs:
  - ./src
```
Run `onefile`, and it will automatically use this configuration.

## Configuration Reference

-   `output` (string): Name of the generated file.
    -   **Default:** `onefile.txt`
-   `input_dirs` (list of strings): Directories to scan. Paths are relative to the config file.
    -   **Default:** `["."]` (the current directory)
-   `exclude_folders` (list of strings): Add extra folders to the built-in exclusion list.
-   `exclude_files` (list of strings): Add specific files to exclude by their full path.

#### Built-in Exclusions
You do not need to manually exclude these:
-   **Folders:** `.git`, `.vscode`, `.idea`, `vendor`, `node_modules`, `bin`, `obj`, `dist`, `build`
-   **Files:** The `onefile` binary, `config-onefile.yml`, the output file itself, `go.mod`, `go.sum`, `package-lock.json`, `yarn.lock`.

## Command-Line Flags

-   `-split`: Splits the final output file into two roughly equal parts (appends `_part1` and `_part2` to the filenames).
    ```bash
    onefile -split
    ```
-   `-config <path>`: Use a configuration file from a non-standard location.
    ```bash
    onefile -config ./configs/custom_config.yml
    ```