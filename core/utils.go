package core

import (
	"os"
	"strings"
)

func ExpandHome(path string) string {
	home, _ := os.UserHomeDir()

	return strings.Replace(path, "~", home, 1)
}
