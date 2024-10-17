package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/kahnwong/workspace-init/core"
)

func main() {
	//repos := core.GetRepos() // [TODO] cross-validate with repos in config

	// config
	config := core.ReadConfig()
	username := config.GitUsername
	workspacePath := ExpandHome(config.WorkspacePath)
	privateKeyFile := ExpandHome(config.PrivateKeyFile)

	_, err := os.Stat(privateKeyFile)
	if err != nil {
		log.Panicf("read file %s failed %s\n", privateKeyFile, err.Error())
	}

	// Clone the given repository to the given directory
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
	if err != nil {
		log.Panicf("generate publickeys failed: %s\n", err.Error())
	}

	// clone repos
	for _, group := range config.Repos {
		for _, repo := range group.Repos {
			repoPath := CreateDir(workspacePath, username, group.Group, repo)
			_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
				Auth:     publicKeys,
				URL:      fmt.Sprintf("git@github.com:%s/%s.git", username, repo),
				Progress: os.Stdout,
			})
			if errors.Is(err, git.ErrRepositoryAlreadyExists) {

			} else if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func CreateDir(workspacePath string, username string, group string, repo string) string {
	repoPath := filepath.Join(workspacePath, username, group, repo)
	err := os.MkdirAll(filepath.Join(repoPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Created directory: %s\n", repoPath)

	return repoPath
}

func ExpandHome(path string) string {
	home, _ := os.UserHomeDir()

	return strings.Replace(path, "~", home, 1)
}
