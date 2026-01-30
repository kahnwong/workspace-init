package core

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cli/go-gh/v2"
)

type RepoList []struct {
	Name string `json:"name"`
}

func getRepos(isArchived bool) ([]string, error) {
	limit := 300
	noArchivedFlag := "--no-archived"
	if isArchived {
		noArchivedFlag = ""
	}
	repoList, _, err := gh.Exec("repo", "list", noArchivedFlag, "--limit", strconv.Itoa(limit), "--json", "name")
	if err != nil {
		return nil, fmt.Errorf("failed to get repo list: %w", err)
	}

	// unmarshal json
	var repoListStruct RepoList
	err = json.Unmarshal(repoList.Bytes(), &repoListStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal repo list: %w", err)
	}

	// create repoListStruct slice
	var repos []string
	for _, repo := range repoListStruct {
		repos = append(repos, repo.Name)
	}

	return repos, nil
}
