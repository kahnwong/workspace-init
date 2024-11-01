package core

import "fmt"

func Terraform() {
	terraformConfig := ""
	categoryConfig := parseCategoryConfig()
	terraformConfig += `
  repos = tomap({
`
	for _, category := range categoryConfig {
		repoArray := fmt.Sprintf("%s = [", category.Group)
		for _, repo := range category.Repos {
			repoArray += fmt.Sprintf(`
	     "%s",`, repo)
		}

		repoArray += `
],
`
		terraformConfig += repoArray
	}

	terraformConfig += `
  })`

	fmt.Println(terraformConfig)
}
