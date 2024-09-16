package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FindGoModRoot searches upwards from the current working directory
// to find the directory containing go.mod.
// findGoModRoot searches upwards from the current working directory
// to find the directory containing go.mod.
func FindGoModRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir { // Reached the root of the file system
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod not found in any parent directory")
}

// DownloadRemoteRepo downloads a ZIP file from GitHub and extracts it to a folder.
// It returns the path to the extracted folder.
func DownloadRemoteRepo(repoPath string) (string, error) {
	repoURL := fmt.Sprintf("https://github.com/%s/archive/refs/heads/main.zip", repoPath)
	zipFileName := strings.ReplaceAll(repoPath, "/", "_") + ".zip"

	// Find the project root where go.mod exists
	projectRoot, err := FindGoModRoot()
	if err != nil {
		return "", fmt.Errorf("failed to find go.mod root: %w", err)
	}

	// Set the download directory in the project root
	downloadDir := filepath.Join(projectRoot, "download")

	// Ensure the download directory exists
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		err := os.Mkdir(downloadDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create download directory: %w", err)
		}
		fmt.Println("Download directory created.")
	}

	zipFilePath := filepath.Join(downloadDir, zipFileName)
	fmt.Printf("Downloading %s to %s...\n", repoURL, zipFilePath)
	err = downloadFile(repoURL, zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}

	extractDir := filepath.Join(downloadDir, strings.TrimSuffix(zipFileName, ".zip"))
	err = extractZip(zipFilePath, extractDir)
	if err != nil {
		return "", fmt.Errorf("failed to extract ZIP file: %w", err)
	}

	fmt.Println("Download and extraction completed successfully.")
	return extractDir, nil
}

// downloadFile downloads a file from a URL and saves it to the specified file path.
func downloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	//defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// extractZip extracts a ZIP file to the specified directory.
func extractZip(zipFilePath, destDir string) error {
	zipFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	for _, file := range zipFile.File {
		fpath := filepath.Join(destDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, file.Mode())
			continue
		}

		outFile, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
