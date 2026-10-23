package core

import (
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestForgeSlugFolder(t *testing.T) {
	for _, slug := range []string{"", "foo"} {
		t.Run("slug="+slug, func(t *testing.T) {
			var org Org
			if err := yaml.Unmarshal([]byte("name: foo\nforge: git.foo.example.com\nforgeSlug: "+slug+"\n"), &org); err != nil {
				t.Fatal(err)
			}
			wantFolder := "foo"
			if slug != "" {
				wantFolder += "-" + slug
			}
			workspace := t.TempDir()
			for _, group := range []string{"", "community"} {
				got, err := createDir(workspace, org.folderName(), group, "repo")
				if err != nil {
					t.Fatal(err)
				}
				if want := filepath.Join(workspace, wantFolder, group, "repo"); got != want {
					t.Errorf("repo path = %q, want %q", got, want)
				}
			}
			if got := repoURL(org.Forge, org.Name, "repo", org.ForgeSSHPort); got != "ssh://git@git.foo.example.com/foo/repo.git" {
				t.Errorf("remote org changed: %q", got)
			}
		})
	}
}

func TestForgeRepoURL(t *testing.T) {
	for _, tt := range []struct {
		name        string
		forgeConfig string
		want        string
	}{
		{"default", "", "git@github.com:foo/repo.git"},
		{"custom", "    forge: git.foo.example.com\n", "ssh://git@git.foo.example.com/foo/repo.git"},
		{"custom port", "    forge: git.foo.example.com:2222\n", "ssh://git@git.foo.example.com:2222/foo/repo.git"},
		{"empty", "    forge: \"\"\n", "git@github.com:foo/repo.git"},
		{"ssh port", "    forge: git.foo.example.com\n    forgeSSHPort: 2222\n", "ssh://git@git.foo.example.com:2222/foo/repo.git"},
		{"default host port", "    forgeSSHPort: 2222\n", "ssh://git@github.com:2222/foo/repo.git"},
		{"override port", "    forge: git.foo.example.com:22\n    forgeSSHPort: 2222\n", "ssh://git@git.foo.example.com:2222/foo/repo.git"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config
			data := "orgs:\n  - name: foo\n    privateKeyFile: ~/.ssh/github\n" + tt.forgeConfig
			if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
				t.Fatal(err)
			}
			org := cfg.Orgs[0]
			if got := repoURL(org.Forge, org.Name, "repo", org.ForgeSSHPort); got != tt.want {
				t.Errorf("repoURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
