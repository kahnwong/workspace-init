package core

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/rs/zerolog/log"
)

func pull(publicKeys *ssh.PublicKeys, path string) (string, error) {
	fmt.Println(path)
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init repo")
	}

	w, err := r.Worktree()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get worktree")
	}

	err = w.Pull(&git.PullOptions{
		Auth:       publicKeys,
		RemoteName: "origin",
	})
	if err != nil {
		fmt.Print(err.Error())
	}

	ref, err := r.Head()
	if err != nil {
		log.Fatal().Err(err).Msg("Error obtaining HEAD")
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Fatal().Err(err).Msg("Error obtaining commit hash")
	}

	return commit.Hash.String(), err
}
