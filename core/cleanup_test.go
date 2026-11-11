package core

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func makeLocalRepo(t *testing.T, path string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(path, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLocalRepositories(t *testing.T) {
	root := t.TempDir()
	a := makeLocalRepo(t, filepath.Join(root, "a"))
	b := makeLocalRepo(t, filepath.Join(root, "group", "b"))
	makeLocalRepo(t, filepath.Join(a, "nested"))
	makeLocalRepo(t, filepath.Join(root, "repo__worktrees"))
	makeLocalRepo(t, filepath.Join(root, "group", "repo__worktrees", "branch"))
	makeLocalRepo(t, filepath.Join(root, "__worktrees", "branch"))
	outside := makeLocalRepo(t, filepath.Join(t.TempDir(), "outside"))
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err != nil {
		t.Fatal(err)
	}
	got, err := localRepositories(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{a, b}) {
		t.Fatalf("got %v, want %v", got, []string{a, b})
	}
	got, err = localRepositories(filepath.Join(root, "missing"))
	if err != nil || len(got) != 0 {
		t.Fatalf("missing directory: %v, %v", got, err)
	}
}

func TestCleanupConfirmation(t *testing.T) {
	for _, answer := range []string{"", "\n", "no\n", "invalid\n", "yes\n", "Y\n"} {
		t.Run(answer, func(t *testing.T) {
			path := makeLocalRepo(t, filepath.Join(t.TempDir(), "repo"))
			var output bytes.Buffer
			if err := confirmCleanup(strings.NewReader(answer), &output, []string{path}, trashDirectory{}); err != nil {
				t.Fatal(err)
			}
			_, err := os.Stat(path)
			wantDeleted := answer == "yes\n" || answer == "Y\n"
			if os.IsNotExist(err) != wantDeleted {
				t.Fatalf("deleted = %v, want %v", os.IsNotExist(err), wantDeleted)
			}
			if !strings.Contains(output.String(), "PERMANENTLY DELETE") {
				t.Fatal("missing permanent deletion warning")
			}
		})
	}
}

func TestCleanupTrash(t *testing.T) {
	trash := trashDirectory{files: t.TempDir(), info: t.TempDir()}
	for range 2 {
		path := makeLocalRepo(t, filepath.Join(t.TempDir(), "repo with spaces"))
		var output bytes.Buffer
		if err := confirmCleanup(strings.NewReader("yes\n"), &output, []string{path}, trash); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("source still exists: %v", err)
		}
	}
	entries, err := os.ReadDir(trash.files)
	if err != nil || len(entries) != 2 {
		t.Fatalf("trash entries: %v, %v", entries, err)
	}
	for _, entry := range entries {
		if _, err := os.Stat(filepath.Join(trash.files, entry.Name(), ".git")); err != nil {
			t.Fatal(err)
		}
		info, err := os.ReadFile(filepath.Join(trash.info, entry.Name()+".trashinfo"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(info), "repo%20with%20spaces\nDeletionDate=") {
			t.Fatalf("invalid trash metadata: %s", info)
		}
	}
}

func TestCleanupTrashFailurePreservesSource(t *testing.T) {
	path := makeLocalRepo(t, filepath.Join(t.TempDir(), "repo"))
	trash := trashDirectory{files: filepath.Join(t.TempDir(), "missing")}
	err := confirmCleanup(strings.NewReader("yes\n"), &bytes.Buffer{}, []string{path}, trash)
	if err == nil {
		t.Fatal("expected trash error")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("source not preserved: %v", err)
	}
}

func TestCleanupInventory(t *testing.T) {
	previousConfig, previousWorkspace := config, workspacePath
	t.Cleanup(func() { config, workspacePath = previousConfig, previousWorkspace })
	workspacePath = t.TempDir()
	config = &Config{Orgs: []Org{{Name: "example"}}}
	root := filepath.Join(workspacePath, "example")
	makeLocalRepo(t, filepath.Join(root, "Active"))
	makeLocalRepo(t, filepath.Join(root, "group", "archived"))
	missing := makeLocalRepo(t, filepath.Join(root, "group", "missing"))
	var output bytes.Buffer
	fetch := func(owner string, includeArchived bool, limit int) ([]string, error) {
		if owner != "example" || !includeArchived || limit != 100000 {
			t.Fatalf("unexpected inventory arguments: %s %v %d", owner, includeArchived, limit)
		}
		return []string{"active", "archived"}, nil
	}
	if err := cleanup(strings.NewReader("no\n"), &output, fetch); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), Blue("example")+":\n- missing\n") || strings.Contains(output.String(), "- Active\n") || strings.Contains(output.String(), "- archived\n") || strings.Contains(output.String(), root) {
		t.Fatalf("unexpected candidates: %s", output.String())
	}
	output.Reset()
	fetch = func(string, bool, int) ([]string, error) { return nil, errors.New("API failure") }
	if err := cleanup(strings.NewReader("yes\n"), &output, fetch); err == nil {
		t.Fatal("expected inventory failure")
	}
	if output.Len() != 0 {
		t.Fatalf("must not prompt on API failure: %s", output.String())
	}
	if _, err := os.Stat(missing); err != nil {
		t.Fatal("API failure removed local repo")
	}
}

func TestCleanupSkipsCustomForge(t *testing.T) {
	previousConfig, previousWorkspace := config, workspacePath
	t.Cleanup(func() { config, workspacePath = previousConfig, previousWorkspace })
	workspacePath = t.TempDir()
	config = &Config{Orgs: []Org{{Name: "example", Forge: "git.example.com"}}}
	path := makeLocalRepo(t, filepath.Join(workspacePath, "example", "repo"))
	t.Setenv("PATH", t.TempDir())
	var output bytes.Buffer
	if err := Cleanup(strings.NewReader("yes\n"), &output); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}
