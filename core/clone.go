package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	Green  = color.New(color.FgHiGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Blue   = color.New(color.FgBlue).SprintFunc()
)

func clone(publicKeys *ssh.PublicKeys, workspacePath string, group string, username string, repo string) {
	repoPath := createDir(workspacePath, username, group, repo)
	repoUrl := fmt.Sprintf("git@github.com:%s/%s.git", username, repo)
	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth: publicKeys,
		URL:  repoUrl,
		//Progress: os.Stdout,
	})
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		fmt.Println(Yellow(fmt.Sprintf("Repo %s already exists", repoUrl)))
		hash, err := pull(publicKeys, repoPath)
		if !errors.Is(err, git.NoErrAlreadyUpToDate) {
		} else {
			fmt.Println(Blue(fmt.Sprintf("Pulled %s", hash)))
		}
	} else if err != nil {
		log.Fatal().Err(err).Msgf("Failed to clone %s", repoUrl)
	} else {
		fmt.Println(Green(fmt.Sprintf("Cloned %s", repoUrl)))
	}
}

func CloneRepos() {
	// config
	username := viper.GetString("gitUsername")
	workspacePath := ExpandHome(viper.GetString("workspacePath"))
	publicKeys := initPublicKey()

	// clone: no category
	noCategoryRepos := viper.GetStringSlice("noCategory")

	var wg sync.WaitGroup
	wg.Add(len(noCategoryRepos))
	for _, repo := range noCategoryRepos {
		go func() {
			clone(publicKeys, workspacePath, "", username, repo)
			wg.Done()
		}()
	}
	wg.Wait()

	// clone: category
	categoryConfig := parseCategoryConfig()
	for _, category := range categoryConfig {
		var wg sync.WaitGroup
		wg.Add(len(category.Repos))
		for _, repo := range category.Repos {
			go func() {
				clone(publicKeys, workspacePath, category.Group, username, repo)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
