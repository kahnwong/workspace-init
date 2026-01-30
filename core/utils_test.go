package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDir(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		workspacePath string
		username      string
		group         string
		repo          string
		expectedPath  string
	}{
		{
			name:          "with group",
			workspacePath: tmpDir,
			username:      "testuser",
			group:         "testgroup",
			repo:          "testrepo",
			expectedPath:  filepath.Join(tmpDir, "testuser", "testgroup", "testrepo"),
		},
		{
			name:          "without group",
			workspacePath: tmpDir,
			username:      "testuser",
			group:         "",
			repo:          "testrepo",
			expectedPath:  filepath.Join(tmpDir, "testuser", "testrepo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createDir(tt.workspacePath, tt.username, tt.group, tt.repo)
			if err != nil {
				t.Fatalf("createDir() error = %v", err)
			}
			if result != tt.expectedPath {
				t.Errorf("createDir() = %q, want %q", result, tt.expectedPath)
			}
			// Verify directory was actually created
			if _, err := os.Stat(result); os.IsNotExist(err) {
				t.Errorf("directory %q was not created", result)
			}
		})
	}
}
