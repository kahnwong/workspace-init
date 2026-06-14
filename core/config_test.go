package core

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfigUnmarshal(t *testing.T) {
	yamlData := `
workspacePath: ~/workspace
orgs:
  - name: testuser
    privateKeyFile: ~/.ssh/id_rsa
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

	if config.WorkspacePath != "~/workspace" {
		t.Errorf("WorkspacePath = %q, want %q", config.WorkspacePath, "~/workspace")
	}
	if len(config.Orgs) != 1 {
		t.Fatalf("Orgs length = %d, want 1", len(config.Orgs))
	}

	org := config.Orgs[0]
	if org.PrivateKeyFile != "~/.ssh/id_rsa" {
		t.Errorf("PrivateKeyFile = %q, want %q", org.PrivateKeyFile, "~/.ssh/id_rsa")
	}
	if org.Name != "testuser" {
		t.Errorf("Name = %q, want %q", org.Name, "testuser")
	}

	// Verify noCategory
	if len(org.NoCategory) != 2 {
		t.Errorf("NoCategory length = %d, want 2", len(org.NoCategory))
	}

	// Verify category
	if len(org.Categories) != 2 {
		t.Errorf("Category length = %d, want 2", len(org.Categories))
	}
	if len(org.Categories) > 0 {
		if org.Categories[0].Group != "group1" {
			t.Errorf("Category[0].Group = %q, want %q", org.Categories[0].Group, "group1")
		}
		if len(org.Categories[0].Repos) != 2 {
			t.Errorf("Category[0].Repos length = %d, want 2", len(org.Categories[0].Repos))
		}
	}

	// Verify excludeRepos
	if len(org.ExcludeRepos) != 1 {
		t.Errorf("ExcludeRepos length = %d, want 1", len(org.ExcludeRepos))
	}
	if len(org.ExcludeRepos) > 0 {
		if org.ExcludeRepos[0].Group != "archived" {
			t.Errorf("ExcludeRepos[0].Group = %q, want %q", org.ExcludeRepos[0].Group, "archived")
		}
		if len(org.ExcludeRepos[0].Repos) != 2 {
			t.Errorf("ExcludeRepos[0].Repos length = %d, want 2", len(org.ExcludeRepos[0].Repos))
		}
	}
}
