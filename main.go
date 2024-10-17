package main

import (
	"fmt"

	"github.com/kahnwong/workspace-init/github"
)

func main() {
	repos := github.GetRepos()
	fmt.Println(repos)
}
