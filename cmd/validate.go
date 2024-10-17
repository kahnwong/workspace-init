package cmd

import (
	"github.com/kahnwong/workspace-init/core"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Display repos not configured in config",
	Run: func(cmd *cobra.Command, args []string) {
		core.Validate()
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
