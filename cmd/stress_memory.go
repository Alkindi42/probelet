package cmd

import (
	"github.com/Alkindi42/probelet/internal/app"
	"github.com/spf13/cobra"
)

type stressMemoryOptions struct {
	duration string
	size     string
}

var stressMemoryOpts stressMemoryOptions

var stressMemoryCmd = &cobra.Command{
	Use:   "memory",
	Args:  cobra.NoArgs,
	Short: "Run a bounded memory stress workload",
	Long: `Run a bounded memory stress workload.

This command is useful for testing autoscaling, resource limits,
and memory-related alerting.

The workload uses the same validation and safety limits as the HTTP API.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		req := app.MemoryStressRequest{
			Size:     stressMemoryOpts.size,
			Duration: stressMemoryOpts.duration,
		}
		result, err := app.RunMemoryStress(cmd.Context(), req)
		if err != nil {
			return err
		}

		cmd.Printf("done: size=%s bytes=%d duration=%s\n", result.Size, result.Bytes, result.Duration)

		return nil
	},
}

func init() {
	stressCmd.AddCommand(stressMemoryCmd)

	stressMemoryCmd.Flags().StringVar(&stressMemoryOpts.size, "size", "", "memory size to allocate")
	stressMemoryCmd.Flags().StringVar(&stressMemoryOpts.duration, "duration", "", "stress duration")

	if err := stressMemoryCmd.MarkFlagRequired("duration"); err != nil {
		panic(err)
	}
	if err := stressMemoryCmd.MarkFlagRequired("size"); err != nil {
		panic(err)
	}
}
