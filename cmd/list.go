package cmd

import (
	"fmt"

	"github.com/shelepuginivan/theme/theme"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available themes",
	RunE:  listCommand,
}

func listCommand(cmd *cobra.Command, _ []string) error {
	prefix, err := cmd.Flags().GetString("prefix")
	if err != nil {
		return err
	}

	t := theme.NewWithPrefix(prefix)

	themes, err := t.Themes()
	if err != nil {
		return err
	}

	for _, t := range themes {
		fmt.Println(t)
	}
	return nil
}
