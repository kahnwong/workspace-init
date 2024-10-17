package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"

	"github.com/kahnwong/workspace-init/core"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "workspace-init",
	Short: "Clone repos into separate folders, depending on grouping",
	Run: func(cmd *cobra.Command, args []string) {
		////repos := core.GetRepos() // [TODO] cross-validate with category in config

		// config
		username := viper.GetString("gitUsername")
		workspacePath := core.ExpandHome(viper.GetString("workspacePath"))
		publicKeys := core.InitPublicKey()

		var category core.Category
		err := viper.UnmarshalKey("category", &category)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to unmarshal category")
		}

		// clone category
		for _, group := range category {
			for _, repo := range group.Repos {
				repoPath := CreateDir(workspacePath, username, group.Group, repo)
				repoUrl := fmt.Sprintf("git@github.com:%s/%s.git", username, repo)
				_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
					Auth:     publicKeys,
					URL:      repoUrl,
					Progress: os.Stdout,
				})
				if errors.Is(err, git.ErrRepositoryAlreadyExists) {

				} else if err != nil {
					log.Fatal().Err(err).Msgf("Failed to clone %s", repoUrl)
				}
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// init
func initConfig() {
	// read config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // [TODO] change to `~/.config/workspace-init`

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

// -- scratch

func CreateDir(workspacePath string, username string, group string, repo string) string {
	repoPath := filepath.Join(workspacePath, username, group, repo)
	err := os.MkdirAll(filepath.Join(repoPath), os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error creating directory %s", repoPath)
	}
	//fmt.Printf("Created directory: %s\n", repoPath)

	return repoPath
}
