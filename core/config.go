package core

import (
	"testing"

	cliBase "github.com/kahnwong/cli-base"
	"github.com/rs/zerolog/log"
)

var config *Config
var workspacePath string

func init() {
	// Unit tests use their own fixtures, not the user's configuration.
	if testing.Testing() {
		return
	}

	var err error
	config, err = cliBase.ReadYaml[Config]("~/.config/workspace-init/config.yaml")
	if err != nil {
		log.Fatal().Msgf("Failed to read config: %v", err)
	}

	workspacePath, err = cliBase.ExpandHome(config.WorkspacePath)
	if err != nil {
		log.Fatal().Msgf("Failed to expand home path: %v", err)
	}
}

type Category struct {
	Group string   `yaml:"group"`
	Repos []string `yaml:"repos"`
}

type ExcludeRepos []struct {
	Group string   `yaml:"group"`
	Repos []string `yaml:"repos"`
}

type Org struct {
	PrivateKeyFile string       `yaml:"privateKeyFile"`
	Name           string       `yaml:"name"`
	NoCategory     []string     `yaml:"noCategory"`
	Categories     []Category   `yaml:"categories"`
	ExcludeRepos   ExcludeRepos `yaml:"excludeRepos"`
}

type Config struct {
	WorkspacePath string `yaml:"workspacePath"`
	Orgs          []Org  `yaml:"orgs"`
}
