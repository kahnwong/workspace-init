package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kahnwong/workspace-init/core"
)

func main() {
	repos := core.GetRepos()
	fmt.Println(repos)

	repoGroups := core.ReadConfig()
	fmt.Println(repoGroups)

	// clone repos
	for _, group := range repoGroups.Repos {
		groupName := group.Group
		CreateDir(groupName)

		//groupRepos := group.Repos
	}
}

func CreateDir(dir string) {
	wd, _ := os.Getwd()                                      // [TODO] remove this
	username := "kahnwong"                                   // [TODO] move to config
	workspacePath := filepath.Join(wd, "Git", username, dir) // [TODO] use real config path

	err := os.MkdirAll(filepath.Join(workspacePath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created directory: %s\n", dir)
}
