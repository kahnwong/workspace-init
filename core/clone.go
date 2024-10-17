package core

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func clone(publicKeys *ssh.PublicKeys, workspacePath string, group string, username string, repo string) {
	repoPath := createDir(workspacePath, username, group, repo)
	repoUrl := fmt.Sprintf("git@github.com:%s/%s.git", username, repo)
	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth:     publicKeys,
		URL:      repoUrl,
		Progress: os.Stdout,
	})
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {

	} else if err != nil {
		log.Fatal().Err(err).Msgf("Failed to clone %s", repoUrl)
	}
}

func CloneRepos() {
	// config
	username := viper.GetString("gitUsername")
	workspacePath := ExpandHome(viper.GetString("workspacePath"))
	publicKeys := initPublicKey()

	var category Category
	err := viper.UnmarshalKey("category", &category)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal category")
	}

	// clone
	for _, group := range category {
		for _, repo := range group.Repos {
			clone(publicKeys, workspacePath, group.Group, username, repo)
		}
	}
}
