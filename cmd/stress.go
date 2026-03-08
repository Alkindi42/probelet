package cmd

import (
	"github.com/spf13/cobra"
)

// stressCmd represents the stress command
var stressCmd = &cobra.Command{
	Use:   "stress",
	Short: "Run controlled system stress workloads",
	Long: `Run controlled system stress workloads such as CPU or memory pressure.

This command is useful for testing autoscaling, resource limits,
alerting, and failure scenarios in containerized or distributed environments.

Stress workloads are bounded and respect Probelet safety limits.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(stressCmd)
}
