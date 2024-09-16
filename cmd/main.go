package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"tkc/llm-repository-loader/internal"
)

func main() {
	remoteRepo := flag.String("remote_repo", "", "Path to the GitHub repository (e.g., tkc/go_bedrock_proxy_server)")
	localRepoPath := flag.String("local_repo_path", "", "Path to the local repo")
	flag.Parse()

	if *remoteRepo == "" && *localRepoPath == "" {
		fmt.Println("Error: You must specify either a remote repository or a local repository.")
		fmt.Println("Usage: go run main.go [-local_repo_path /path/to/local/repo] [-remote_repo tkc/go_bedrock_proxy_server]")
		os.Exit(1)
	}

	// Find the project root where go.mod exists
	projectRoot, err := internal.FindGoModRoot()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Define the output directory in the same location as go.mod
	outputDir := filepath.Join(projectRoot, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	var repoPath string
	if *remoteRepo != "" {
		// Download and extract remote repo if specified
		repoPath, err = internal.DownloadRemoteRepo(*remoteRepo)
		if err != nil {
			fmt.Printf("Error downloading repository: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Repository downloaded and extracted to %s\n", repoPath)
	} else if *localRepoPath != "" {
		// Use local repository if specified
		repoPath = *localRepoPath
		fmt.Printf("Using local repository at %s\n", repoPath)
	}

	ignoreFilePath := filepath.Join(projectRoot, ".loader_ignores")

	var ignoreList []string
	if _, err := os.Stat(ignoreFilePath); err == nil {
		ignoreList = internal.GetIgnoreList(ignoreFilePath)
	}

	fmt.Print(ignoreList)

	outputFilePath := filepath.Join(outputDir, strings.ReplaceAll(filepath.Base(repoPath), "/", "_")+".txt")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}

	internal.ProcessRepository(repoPath, ignoreList, outputFile)

	outputFile.WriteString("--END--\n")
	fmt.Printf("Repository contents written to %s\n", outputFilePath)
}
