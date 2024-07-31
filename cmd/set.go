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
	RunE:  setCommand,
}

func setCommand(cmd *cobra.Command, args []string) error {
	prefix, err := cmd.Flags().GetString("prefix")
	if err != nil {
		return err
	}

	t := theme.NewWithPrefix(prefix)

	for _, err := range t.Set(args[0]) {
		cmd.PrintErrln(err)
	}
	return nil
}
