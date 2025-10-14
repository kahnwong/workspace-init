package core

import (
	"fmt"
)

func Validate() {
	reposActive := getRepos(false)
	excludeRepos := config.ExcludeRepos
	reposExcluded := subtractArrays(reposActive, excludeRepos)

	var reposConfig []string
	reposConfig = append(reposConfig, config.NoCategory...)
	for _, category := range config.Category {
		reposConfig = append(reposConfig, category.Repos...)
	}

	reposNotInConfig := subtractArrays(reposExcluded, reposConfig)
	if len(reposNotInConfig) > 0 {
		fmt.Println("Following repos are not in config:")

		for _, repo := range reposNotInConfig {
			fmt.Printf("- %s\n", repo)
		}
	} else {
		fmt.Println("All repos are declared in config")
	}

}

func subtractArrays(array1, array2 []string) []string {
	var result []string

	// Create a map to efficiently check if an element is in array2
	exists := make(map[string]bool)
	for _, elem := range array2 {
		exists[elem] = true
	}

	// Iterate through array1 and add elements not in array2 to the result
	for _, elem := range array1 {
		if !exists[elem] {
			result = append(result, elem)
		}
	}

	return result
}
