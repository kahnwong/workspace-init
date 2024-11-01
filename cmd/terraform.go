package cmd

import (
	"github.com/kahnwong/workspace-init/core"
	"github.com/spf13/cobra"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Display repos config to be used for forgejo mirrors in terraform",
	Run: func(cmd *cobra.Command, args []string) {
		core.Terraform()
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)
}
