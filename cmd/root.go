package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringP("prefix", "p", filepath.Join(xdg.ConfigHome, "theme"), "Directory where themes are stored")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress warnings and subprocess output")
}

var rootCmd = &cobra.Command{
	Use:   "theme",
	Short: "A very simple theme switcher",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
