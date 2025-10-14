package core

import (
	"errors"
	"fmt"
	"sync"

	cliBase "github.com/kahnwong/cli-base"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/rs/zerolog/log"
)

var (
	Green  = color.New(color.FgHiGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
)

func clone(publicKeys *ssh.PublicKeys, workspacePath string, group string, repo string, username string) {
	if group == "" { // noCategory
		group = "."
	}
	repoPath := createDir(workspacePath, username, group, repo)
	repoUrl := fmt.Sprintf("git@github.com:%s/%s.git", username, repo)
	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth: publicKeys,
		URL:  repoUrl,
		//Progress: os.Stdout,
	})
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		fmt.Println(Yellow(fmt.Sprintf("Repo %s already exists", repoUrl)))
	} else if err != nil {
		log.Fatal().Msgf("Failed to clone %s", repoUrl)
	} else {
		fmt.Println(Green(fmt.Sprintf("Cloned %s", repoUrl)))
	}
}

func CloneRepos() {
	// config
	username := config.GitUsername
	workspacePath := cliBase.ExpandHome(config.WorkspacePath)
	publicKeys := initPublicKey()

	// clone
	noCategoryConfig := config.NoCategory
	var wg sync.WaitGroup
	wg.Add(len(noCategoryConfig))
	for _, repo := range noCategoryConfig {
		go func() {
			clone(publicKeys, workspacePath, "", repo, username)
			wg.Done()
		}()
	}
	wg.Wait()

	categoryConfig := config.Category
	for _, category := range categoryConfig {
		var wg sync.WaitGroup
		wg.Add(len(category.Repos))
		for _, repo := range category.Repos {
			go func() {
				clone(publicKeys, workspacePath, category.Group, repo, username)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
