package cmd

import (
	"github.com/shelepuginivan/theme/theme"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(randomCmd)
}

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Set a random available theme",
	Run:   randomCommand,
}

func randomCommand(cmd *cobra.Command, _ []string) {
	var c theme.Config

	c.Prefix, _ = cmd.Flags().GetString("prefix")
	c.Quiet, _ = cmd.Flags().GetBool("quiet")

	theme.NewWithConfig(c).Random()
}
