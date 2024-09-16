package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetIgnoreList(ignoreFilePath string) []string {
	var ignoreList []string

	file, err := os.Open(ignoreFilePath)
	if err != nil {
		return ignoreList
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if runtime.GOOS == "windows" {
			line = strings.ReplaceAll(line, "/", "\\")
		}
		ignoreList = append(ignoreList, strings.TrimSpace(line))
	}

	return ignoreList
}

// IsGitSubfolder checks if a given path is within a .git/ directory
func IsGitSubfolder(pattern, path string) bool {
	// Convert the path to its absolute form to ensure consistent comparison
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Check if .git/ is part of the path
	gitDir := filepath.Join(pattern, string(filepath.Separator))

	// Check if the path contains ".git/" as part of its directory
	return strings.Contains(absPath, gitDir)
}

// ShouldIgnore checks if a file path matches any pattern in the ignore list.
func ShouldIgnore(filePath string, ignoreList []string) bool {
	for _, pattern := range ignoreList {
		fmt.Println("/")
		fmt.Printf(filePath)

		if IsGitSubfolder(pattern, filePath) {
			return true
		}

		match, _ := filepath.Match(pattern, filePath)
		if match {
			return true
		}
	}
	return false
}
