package core

import (
	"errors"
	"fmt"
	"net"
	"strconv"
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

func repoURL(forge, username, repo string, port int) string {
	if port != 0 {
		if forge == "" {
			forge = "github.com"
		}
		if host, _, err := net.SplitHostPort(forge); err == nil {
			forge = host
		}
		forge = net.JoinHostPort(forge, strconv.Itoa(port))
	}
	if forge != "" {
		return fmt.Sprintf("ssh://git@%s/%s/%s.git", forge, username, repo)
	}
	return fmt.Sprintf("git@github.com:%s/%s.git", username, repo)
}

func (org Org) folderName() string {
	if org.ForgeSlug == "" {
		return org.Name
	}
	return org.Name + "-" + org.ForgeSlug
}

func clone(publicKeys *ssh.PublicKeys, workspacePath string, group string, repo string, org Org) error {
	if group == "" { // noCategory
		group = "."
	}
	repoPath, err := createDir(workspacePath, org.folderName(), group, repo)
	if err != nil {
		return err
	}
	repoUrl := repoURL(org.Forge, org.Name, repo, org.ForgeSSHPort)
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
	for _, org := range config.Orgs {
		if err := cloneOrg(org); err != nil {
			return err
		}
	}

	return nil
}

func cloneOrg(org Org) error {
	publicKeys, err := initPublicKey(org)
	if err != nil {
		return err
	}

	// clone
	noCategoryConfig := org.NoCategory
	var wg sync.WaitGroup
	errChan := make(chan error, len(noCategoryConfig))
	wg.Add(len(noCategoryConfig))
	for _, repo := range noCategoryConfig {
		go func(r string) {
			defer wg.Done()
			if cloneErr := clone(publicKeys, workspacePath, "", r, org); cloneErr != nil {
				errChan <- cloneErr
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

	categoryConfig := org.Categories
	for _, category := range categoryConfig {
		errChan = make(chan error, len(category.Repos))
		wg.Add(len(category.Repos))
		for _, repo := range category.Repos {
			go func(r string, g string) {
				defer wg.Done()
				if cloneErr := clone(publicKeys, workspacePath, g, r, org); cloneErr != nil {
					errChan <- cloneErr
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
