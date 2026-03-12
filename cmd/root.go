package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the probelet command.
var rootCmd = &cobra.Command{
	Use:   "probelet",
	Short: "HTTP service and CLI for failure simulation and system stress testing",
	Long: `Probelet is a lightweight tool for simulating application behavior,
failure scenarios, and controlled system stress workloads.

It can run as an HTTP service or as a local CLI for bounded CPU, memory,
and disk stress operations.`,
	SilenceErrors: true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
