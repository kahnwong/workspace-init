package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/kahnwong/workspace-init/github"
)

type RepoGroups struct {
	Repos []struct {
		Group string   `json:"group"`
		Repos []string `json:"repos"`
	} `json:"repos"`
}

func main() {
	repos := github.GetRepos()
	fmt.Println(repos)

	// read yaml
	configFile, err := os.ReadFile("repos.yaml")
	if err != nil {
		panic(err)
	}

	// Unmarshal the YAML data into a Config struct
	var repoGroups RepoGroups
	err = yaml.Unmarshal(configFile, &repoGroups)
	if err != nil {
		panic(err)
	}

	fmt.Println(repoGroups)
}
