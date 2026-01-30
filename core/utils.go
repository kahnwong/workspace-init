package core

import (
	"fmt"
	"os"
	"path/filepath"
)

func createDir(workspacePath string, username string, group string, repo string) (string, error) {
	repoPath := ""
	if group != "" {
		repoPath = filepath.Join(workspacePath, username, group, repo)
	} else {
		repoPath = filepath.Join(workspacePath, username, repo)
	}

	err := os.MkdirAll(filepath.Join(repoPath), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating directory %s: %w", repoPath, err)
	}

	return repoPath, nil
}
