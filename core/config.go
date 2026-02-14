package core

import (
	"flag"

	cliBase "github.com/kahnwong/cli-base"
	"github.com/rs/zerolog/log"
)

var config *Config
var workspacePath string

func init() {
	var err error
	config, err = cliBase.ReadYaml[Config]("~/.config/workspace-init/config.yaml")
	if err != nil {
		if flag.Lookup("test.v") != nil {
			log.Warn().Msgf("Failed to read config (test mode): %v", err)
			return
		}
		log.Fatal().Msgf("Failed to read config: %v", err)
	}

	workspacePath, err = cliBase.ExpandHome(config.WorkspacePath)
	if err != nil {
		if flag.Lookup("test.v") != nil {
			log.Warn().Msgf("Failed to expand home path (test mode): %v", err)
			return
		}
		log.Fatal().Msgf("Failed to expand home path: %v", err)
	}
}

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
