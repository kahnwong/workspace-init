package core

import (
	cliBase "github.com/kahnwong/cli-base"
)

var config = cliBase.ReadYaml[Config]("~/.config/workspace-init/config.yaml")

type Category []struct {
	Group string   `yaml:"group"`
	Repos []string `yaml:"repos"`
}

type ExcludeRepos []struct {
	Group string   `yaml:"group"`
	Repos []string `yaml:"repos"`
}

type Config struct {
	PrivateKeyFile string       `yaml:"privateKeyFile"`
	WorkspacePath  string       `yaml:"workspacePath"`
	GitUsername    string       `yaml:"gitUsername"`
	NoCategory     []string     `yaml:"noCategory"`
	Category       Category     `yaml:"category"`
	ExcludeRepos   ExcludeRepos `yaml:"excludeRepos"`
}
