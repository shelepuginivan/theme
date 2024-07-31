package cmd

import (
	"github.com/shelepuginivan/theme/theme"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available themes",
	Run:   listCommand,
}

func listCommand(cmd *cobra.Command, _ []string) {
	var c theme.Config

	c.Prefix, _ = cmd.Flags().GetString("prefix")
	c.Quiet, _ = cmd.Flags().GetBool("quiet")

	theme.NewWithConfig(c).List()
}
