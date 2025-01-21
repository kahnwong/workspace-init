package cmd

import (
	"os"

	"github.com/kahnwong/workspace-init/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "workspace-init",
	Short: "Clone repos into separate folders, depending on grouping",
	Run: func(cmd *cobra.Command, args []string) {
		core.CloneRepos()
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

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
