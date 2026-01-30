package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var (
	Green  = color.New(color.FgHiGreen).SprintFunc()
	Blue   = color.New(color.FgBlue).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
)

func clone(publicKeys *ssh.PublicKeys, workspacePath string, group string, repo string, username string) error {
	if group == "" { // noCategory
		group = "."
	}
	repoPath, err := createDir(workspacePath, username, group, repo)
	if err != nil {
		return err
	}
	repoUrl := fmt.Sprintf("git@github.com:%s/%s.git", username, repo)
	_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth: publicKeys,
		URL:  repoUrl,
		//Progress: os.Stdout,
	})
	if errors.Is(err, git.ErrRepositoryAlreadyExists) { // repo already exists, performing git fetch
		var r *git.Repository
		// Open existing repository and fetch from origin
		r, err = git.PlainOpen(repoPath)
		if err != nil {
			return fmt.Errorf("failed to open existing repo %s: %w", repoPath, err)
		}

		err = r.Fetch(&git.FetchOptions{
			Auth:       publicKeys,
			RemoteName: "origin",
		})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return fmt.Errorf("failed to fetch origin for %s: %w", repoUrl, err)
		} else {
			fmt.Println(Blue(fmt.Sprintf("Fetched origin for %s", repoUrl)))
		}
	} else if err != nil {
		return fmt.Errorf("failed to clone %s: %w", repoUrl, err)
	} else { // repo does not exist locally, perform git clone
		fmt.Println(Green(fmt.Sprintf("Cloned %s", repoUrl)))
	}
	return nil
}

func CloneRepos() error {
	// config
	username := config.GitUsername
	publicKeys, err := initPublicKey()
	if err != nil {
		return err
	}

	// clone
	noCategoryConfig := config.NoCategory
	var wg sync.WaitGroup
	errChan := make(chan error, len(noCategoryConfig))
	wg.Add(len(noCategoryConfig))
	for _, repo := range noCategoryConfig {
		go func(r string) {
			defer wg.Done()
			if err = clone(publicKeys, workspacePath, "", r, username); err != nil {
				errChan <- err
			}
		}(repo)
	}
	wg.Wait()
	close(errChan)

	for err = range errChan {
		if err != nil {
			return err
		}
	}

	categoryConfig := config.Category
	for _, category := range categoryConfig {
		errChan = make(chan error, len(category.Repos))
		wg.Add(len(category.Repos))
		for _, repo := range category.Repos {
			go func(r string, g string) {
				defer wg.Done()
				if err = clone(publicKeys, workspacePath, g, r, username); err != nil {
					errChan <- err
				}
			}(repo, category.Group)
		}
		wg.Wait()
		close(errChan)

		for err = range errChan {
			if err != nil {
				return err
			}
		}
	}

	return nil
}
