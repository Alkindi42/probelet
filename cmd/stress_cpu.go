package cmd

import (
	"github.com/Alkindi42/probelet/internal/app"
	"github.com/spf13/cobra"
)

type stressCPUOptions struct {
	duration string
	cores    string
}

var stressCPUOpts stressCPUOptions

var stressCPUCmd = &cobra.Command{
	Use:   "cpu",
	Args:  cobra.NoArgs,
	Short: "Run CPU stress workload",
	Long: `Run a bounded CPU stress workload.

This command is useful for testing autoscaling, throttling,
resource limits, and CPU-related alerting.

The workload uses the same validation and safety limits as the HTTP API.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		req := app.CPUStressRequest{
			Cores:    stressCPUOpts.cores,
			Duration: stressCPUOpts.duration,
		}
		result, err := app.RunCPUStress(cmd.Context(), req)
		if err != nil {
			return err
		}

		cmd.Printf("done: cores=%d duration=%s\n", result.Cores, result.Duration)

		return nil
	},
}

func init() {
	stressCmd.AddCommand(stressCPUCmd)

	stressCPUCmd.Flags().StringVar(&stressCPUOpts.cores, "cores", "", "number of CPU cores or 'max'")
	stressCPUCmd.Flags().StringVar(&stressCPUOpts.duration, "duration", "", "stress duration")

	if err := stressCPUCmd.MarkFlagRequired("duration"); err != nil {
		panic(err)
	}
}
