package core

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RepoGroups struct {
	Repos []struct {
		Group string   `json:"group"`
		Repos []string `json:"repos"`
	} `json:"repos"`
}

func ReadConfig() RepoGroups {
	configFile, err := os.ReadFile("repos.yaml")
	if err != nil {
		panic(err)
	}

	var repoGroups RepoGroups
	err = yaml.Unmarshal(configFile, &repoGroups)
	if err != nil {
		panic(err)
	}

	return repoGroups
}
