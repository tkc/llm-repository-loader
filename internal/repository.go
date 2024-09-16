package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ProcessRepository(repoPath string, ignoreList []string, outputFile *os.File) {
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relativeFilePath, err := filepath.Rel(repoPath, path)
			if err != nil {
				return err
			}

			if !ShouldIgnore(relativeFilePath, ignoreList) {
				contents, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				_, err = outputFile.WriteString("----\n")
				if err != nil {
					return err
				}
				_, err = outputFile.WriteString(fmt.Sprintf("%s\n", relativeFilePath))
				if err != nil {
					return err
				}
				_, err = outputFile.WriteString(fmt.Sprintf("%s\n", contents))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error processing repository: %v\n", err)
	}
}
