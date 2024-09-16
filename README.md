# LLM Repository Loader

LLM Repository Loader is a command-line tool written in Go that reads a Git repository, applies an ignore list, and outputs the repository's file structure and contents to a specified text file. The tool allows you to either specify a local repository path or download a GitHub repository directly, then process the files in that repository. Files are ignored based on patterns specified in an ignore file (.loader_ignores), and the final output is written to a text file.

## Features
- Local and Remote Repositories: You can provide either a local path to a Git repository or download one directly from GitHub.
- Custom Ignore List: You can specify files and directories to ignore by providing patterns in a .loader_ignores file.
- File Processing: The tool processes each file in the repository and writes the contents to a specified output file, including file paths.
- Support for Preamble: You can provide a preamble text file to be included at the beginning of the output.
- Cross-platform: Supports both Windows and Unix-based systems.

## Usage

Command-line Options

- local_repo_path (Optional): The path to a local Git repository.
- remote_repo (Optional): The path to a GitHub repository in the format user/repo_name. The repository will be downloaded and processed.

## Examples

Processing a local repository:
```bash
go run main.go -local_repo_path /path/to/local/repository
```

Downloading and processing a GitHub repository:

```bash
go run main.go -remote_repo tkc/go_bedrock_proxy_server
```

## .loader_ignores File
The .loader_ignores file specifies patterns for files or directories that should be excluded from processing. The patterns follow standard glob-style matching.

```lua
*.log
output/
.git/
.idea/
```

## Output

```txt
----
<relative-file-path>
<file-contents>
----
```

If a preamble is provided, it will be included at the beginning of the output file.

