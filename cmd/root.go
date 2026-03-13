package cmd

import (
	"fmt"
	"os"
	"strings"

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
	SilenceUsage:  true,
}

func Execute() {
	cmd, err := rootCmd.ExecuteC()
	if err == nil {
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, "Error:", err)

	msg := err.Error()
	if strings.Contains(msg, "unknown command") || strings.Contains(msg, "unknown flag") {
		_, _ = fmt.Fprintln(os.Stderr)
		_ = cmd.Help()
	}

	os.Exit(1)
}
