package github

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/cli/go-gh/v2"
)

type RepoList []struct {
	Name string `json:"name"`
}

func GetRepos() []string {
	limit := 5 // [TODO] set to 300 later
	repoList, _, err := gh.Exec("repo", "list", "--no-archived", "--limit", strconv.Itoa(limit), "--json", "name")
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal json
	var repoListStruct RepoList
	err = json.Unmarshal(repoList.Bytes(), &repoListStruct)
	if err != nil {
		log.Fatal(err)
	}

	// create repoListStruct slice
	var repos []string
	for _, repo := range repoListStruct {
		repos = append(repos, repo.Name)
	}

	return repos
}
