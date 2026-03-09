package cmd

import (
	"github.com/Alkindi42/probelet/internal/version"
	"github.com/spf13/cobra"
)

// versionOptions holds the flags for the version command.
type versionOptions struct {
	short bool
}

var versionOpts versionOptions

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build version information",
	Long:  "Print build version information for the Probelet binary.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if versionOpts.short {
			cmd.Println(version.Version)
			return nil
		}

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
	versionCmd.Flags().BoolVar(&versionOpts.short, "short", false, "Print only the version number")

	rootCmd.AddCommand(versionCmd)
}
