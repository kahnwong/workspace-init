package core

import (
	"encoding/json"
	"strconv"

	"github.com/cli/go-gh/v2"
	"github.com/rs/zerolog/log"
)

type RepoList []struct {
	Name string `json:"name"`
}

func getRepos() []string {
	limit := 300
	repoList, _, err := gh.Exec("repo", "list", "--no-archived", "--limit", strconv.Itoa(limit), "--json", "name")
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting repo list")
	}

	// unmarshal json
	var repoListStruct RepoList
	err = json.Unmarshal(repoList.Bytes(), &repoListStruct)
	if err != nil {
		log.Fatal().Err(err).Msg("Error unmarshalling repo list")
	}

	// create repoListStruct slice
	var repos []string
	for _, repo := range repoListStruct {
		repos = append(repos, repo.Name)
	}

	return repos
}
