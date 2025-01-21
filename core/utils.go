package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func ExpandHome(path string) string {
	home, _ := os.UserHomeDir()

	return strings.Replace(path, "~", home, 1)
}

func createDir(workspacePath string, username string, group string, repo string) string {
	repoPath := ""
	if group != "" {
		repoPath = filepath.Join(workspacePath, username, group, repo)
	} else {
		repoPath = filepath.Join(workspacePath, username, repo)
	}

	err := os.MkdirAll(filepath.Join(repoPath), os.ModePerm)
	if err != nil {
		log.Fatal().Msgf("Error creating directory %s", repoPath)
	}

	return repoPath
}
