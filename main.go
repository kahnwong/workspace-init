package main

import (
	"fmt"

	"github.com/kahnwong/workspace-init/core"
)

func main() {
	repos := core.GetRepos()
	fmt.Println(repos)

	repoGroups := core.ReadConfig()
	fmt.Println(repoGroups)
}
