package core

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfigUnmarshal(t *testing.T) {
	yamlData := `
privateKeyFile: ~/.ssh/id_rsa
workspacePath: ~/workspace
gitUsername: testuser
noCategory:
  - repo1
  - repo2
category:
  - group: group1
    repos:
      - repo3
      - repo4
  - group: group2
    repos:
      - repo5
excludeRepos:
  - group: archived
    repos:
      - oldrepo1
      - oldrepo2
`

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify basic fields
	if config.PrivateKeyFile != "~/.ssh/id_rsa" {
		t.Errorf("PrivateKeyFile = %q, want %q", config.PrivateKeyFile, "~/.ssh/id_rsa")
	}
	if config.WorkspacePath != "~/workspace" {
		t.Errorf("WorkspacePath = %q, want %q", config.WorkspacePath, "~/workspace")
	}
	if config.GitUsername != "testuser" {
		t.Errorf("GitUsername = %q, want %q", config.GitUsername, "testuser")
	}

	// Verify noCategory
	if len(config.NoCategory) != 2 {
		t.Errorf("NoCategory length = %d, want 2", len(config.NoCategory))
	}

	// Verify category
	if len(config.Category) != 2 {
		t.Errorf("Category length = %d, want 2", len(config.Category))
	}
	if len(config.Category) > 0 {
		if config.Category[0].Group != "group1" {
			t.Errorf("Category[0].Group = %q, want %q", config.Category[0].Group, "group1")
		}
		if len(config.Category[0].Repos) != 2 {
			t.Errorf("Category[0].Repos length = %d, want 2", len(config.Category[0].Repos))
		}
	}

	// Verify excludeRepos
	if len(config.ExcludeRepos) != 1 {
		t.Errorf("ExcludeRepos length = %d, want 1", len(config.ExcludeRepos))
	}
	if len(config.ExcludeRepos) > 0 {
		if config.ExcludeRepos[0].Group != "archived" {
			t.Errorf("ExcludeRepos[0].Group = %q, want %q", config.ExcludeRepos[0].Group, "archived")
		}
		if len(config.ExcludeRepos[0].Repos) != 2 {
			t.Errorf("ExcludeRepos[0].Repos length = %d, want 2", len(config.ExcludeRepos[0].Repos))
		}
	}
}
