package core

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//type Config struct {
//	PrivateKeyFile string   `yaml:"privateKeyFile"`
//	WorkspacePath  string   `yaml:"workspacePath"`
//	GitUsername    string   `yaml:"gitUsername"`
//	Category       Category `yaml:"repos"`
//}

type Category []struct {
	Group string   `yaml:"group"`
	Repos []string `yaml:"repos"`
}

func parseCategoryConfig() Category {
	var category Category
	err := viper.UnmarshalKey("category", &category)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal category")
	}

	return category
}
