package cmd

import (
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/shelepuginivan/theme/theme"
	"github.com/spf13/cobra"
)

func init() {
	setCmd.Flags().BoolP("quiet", "q", false, "Suppress warnings and subprocess output")
	setCmd.Flags().StringP("prefix", "p", filepath.Join(xdg.ConfigHome, "theme"), "Directory where themes are stored")

	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set <theme>",
	Short: "Set a theme",
	Args:  cobra.ExactArgs(1),
	Run:   setCommand,
}

func setCommand(cmd *cobra.Command, args []string) {
	var c theme.Config

	c.Prefix, _ = cmd.Flags().GetString("prefix")
	c.Quiet, _ = cmd.Flags().GetBool("quiet")

	theme.NewWithConfig(c).Set(args[0])
}
