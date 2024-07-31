package cmd

import (
	"github.com/shelepuginivan/theme/theme"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set <theme>",
	Short: "Set a theme",
	Args:  cobra.ExactArgs(1),
	Run:   setCommand,
}

func setCommand(_ *cobra.Command, args []string) {
	theme.New().Set(args[0])
}
