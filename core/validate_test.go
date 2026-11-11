package core

import "testing"

func TestValidateOrgSkipsCustomForge(t *testing.T) {
	// No GitHub CLI should be invoked for repositories on another forge.
	t.Setenv("PATH", t.TempDir())

	org := Org{Name: "example", Forge: "git.example.com"}
	if err := validateOrg(org); err != nil {
		t.Fatalf("validateOrg() = %v, want nil for a custom forge", err)
	}
}
