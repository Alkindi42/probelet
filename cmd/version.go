package cmd

import (
	"github.com/Alkindi42/probelet/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build version information",
	Long:  "Print build version information for the Probelet binary.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Printf(
			"version: %s\ncommit: %s\nbuild_date: %s\n",
			version.Version,
			version.Commit,
			version.BuildDate,
		)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
