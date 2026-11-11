package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"
)

// Cleanup collects all candidates before prompting or modifying any files.
func Cleanup(in io.Reader, out io.Writer) error {
	return cleanup(in, out, getReposWithLimit)
}

func cleanup(in io.Reader, out io.Writer, fetchRepos func(string, bool, int) ([]string, error)) error {
	var candidates []string
	owners := make(map[string]string)
	for _, org := range config.Orgs {
		if org.Forge != "" {
			continue
		}
		root, err := filepath.Abs(filepath.Join(workspacePath, org.folderName()))
		if err != nil {
			return err
		}
		base, err := filepath.Abs(workspacePath)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(base, root)
		if err != nil || rel == "." || !filepath.IsLocal(rel) {
			return fmt.Errorf("unsafe organization directory: %s", root)
		}
		local, err := localRepositories(root)
		if err != nil {
			return err
		}
		if len(local) == 0 {
			continue
		}
		// gh repo list paginates up to this limit; reaching it fails closed.
		remote, err := fetchRepos(org.Name, true, 100000)
		if err != nil {
			return err
		}
		known := make(map[string]bool, len(remote))
		for _, name := range remote {
			known[strings.ToLower(name)] = true
		}
		for _, path := range local {
			if !known[strings.ToLower(filepath.Base(path))] {
				candidates = append(candidates, path)
				owners[path] = org.Name
			}
		}
	}
	slices.Sort(candidates)
	candidates = slices.Compact(candidates)
	if len(candidates) > 0 {
		fmt.Fprintln(out, "Local repositories not visible on GitHub:")
		groups := make(map[string][]string)
		for _, path := range candidates {
			groups[owners[path]] = append(groups[owners[path]], filepath.Base(path))
		}
		for _, org := range slices.Sorted(maps.Keys(groups)) {
			fmt.Fprintf(out, "%s:\n", Blue(org))
			slices.Sort(groups[org])
			for _, repo := range groups[org] {
				fmt.Fprintf(out, "- %s\n", repo)
			}
		}
	}
	return confirmCleanup(in, out, candidates, availableTrash())
}

// Only actual repository directories are candidates; never follow symlinks or
// descend into a repository (which might contain nested repos or submodules).
func localRepositories(root string) ([]string, error) {
	var repos []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			if path == root && errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return err
		}
		if !entry.IsDir() {
			return nil
		}
		if strings.HasSuffix(entry.Name(), "__worktrees") {
			return filepath.SkipDir
		}
		_, err = os.Lstat(filepath.Join(path, ".git"))
		if err == nil {
			if path != root {
				repos = append(repos, path)
			}
			return filepath.SkipDir
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return nil
	})
	return repos, err
}

type trashDirectory struct {
	files string
	info  string
}

func availableTrash() trashDirectory {
	home, err := os.UserHomeDir()
	if err != nil {
		return trashDirectory{}
	}
	var trash trashDirectory
	switch runtime.GOOS {
	case "darwin":
		trash.files = filepath.Join(home, ".Trash")
	case "linux":
		data := os.Getenv("XDG_DATA_HOME")
		if !filepath.IsAbs(data) {
			data = filepath.Join(home, ".local", "share")
		}
		trash.files = filepath.Join(data, "Trash", "files")
		trash.info = filepath.Join(data, "Trash", "info")
	default:
		return trashDirectory{}
	}
	for _, path := range []string{trash.files, trash.info} {
		if path == "" {
			continue
		}
		stat, err := os.Stat(path)
		if err != nil || !stat.IsDir() {
			return trashDirectory{}
		}
	}
	return trash
}

func confirmCleanup(in io.Reader, out io.Writer, paths []string, trash trashDirectory) error {
	if len(paths) == 0 {
		fmt.Fprintln(out, "No local-only repositories found.")
		return nil
	}
	if trash.files != "" {
		fmt.Fprint(out, Yellow(fmt.Sprintf("Move these repositories to trash (%s)? [y/N] ", trash.files)))
	} else {
		fmt.Fprint(out, "Trash unavailable. PERMANENTLY DELETE these repositories, including uncommitted files? [y/N] ")
	}
	scanner := bufio.NewScanner(in)
	if !scanner.Scan() {
		return scanner.Err()
	}
	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if answer != "y" && answer != "yes" {
		fmt.Fprintln(out, "Cancelled.")
		return nil
	}
	for _, path := range paths {
		var err error
		if trash.files != "" {
			err = moveToTrash(path, trash)
		} else {
			err = os.RemoveAll(path)
		}
		if err != nil {
			return fmt.Errorf("cleanup %s: %w", path, err)
		}
		fmt.Fprintf(out, "Removed %s\n", path)
	}
	return nil
}

func moveToTrash(path string, trash trashDirectory) error {
	// Reserve a unique destination, so identically named repos never overwrite.
	dest, err := os.MkdirTemp(trash.files, filepath.Base(path)+"-")
	if err != nil {
		return err
	}
	infoPath := ""
	if trash.info != "" {
		infoPath = filepath.Join(trash.info, filepath.Base(dest)+".trashinfo")
		info, err := os.OpenFile(infoPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			os.Remove(dest)
			return err
		}
		_, writeErr := fmt.Fprintf(info, "[Trash Info]\nPath=%s\nDeletionDate=%s\n", (&url.URL{Path: path}).EscapedPath(), time.Now().Format("2006-01-02T15:04:05"))
		closeErr := info.Close()
		if err := errors.Join(writeErr, closeErr); err != nil {
			os.Remove(infoPath)
			os.Remove(dest)
			return err
		}
	}
	// Go's Rename requires a non-existing destination for directories.
	if err := os.Remove(dest); err != nil {
		if infoPath != "" {
			os.Remove(infoPath)
		}
		return err
	}
	// A failed move (including across filesystems) must not trigger hard deletion.
	if err := os.Rename(path, dest); err != nil {
		if infoPath != "" {
			os.Remove(infoPath)
		}
		return err
	}
	return nil
}
