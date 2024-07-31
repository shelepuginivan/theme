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
	RunE:  randomCommand,
}

func randomCommand(cmd *cobra.Command, _ []string) error {
	prefix, err := cmd.Flags().GetString("prefix")
	if err != nil {
		return err
	}

	t := theme.NewWithPrefix(prefix)

	r, err := t.Random()
	if err != nil {
		return err
	}

	for _, err := range t.Set(r) {
		cmd.PrintErrln(err)
	}
	return nil
}
