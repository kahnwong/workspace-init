package core

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	PrivateKeyFile string `yaml:"privateKeyFile"`
	WorkspacePath  string `yaml:"workspacePath"`
	GitUsername    string `yaml:"gitUsername"`
	Repos          []struct {
		Group string   `yaml:"group"`
		Repos []string `yaml:"repos"`
	} `yaml:"repos"`
}

func ReadConfig() Config {
	configFile, err := os.ReadFile("workspace-init.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
