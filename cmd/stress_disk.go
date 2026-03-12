package cmd

import (
	"github.com/Alkindi42/probelet/internal/app"
	"github.com/spf13/cobra"
)

type stressDiskOptions struct {
	duration string
	size     string
}

var stressDiskOpts stressDiskOptions

var stressDiskCmd = &cobra.Command{
	Use:   "disk",
	Args:  cobra.NoArgs,
	Short: "Run a bounded disk stress workload",
	Long: `Run a bounded disk stress workload.

This command is useful for testing disk pressure, ephemeral storage limits,
eviction behavior, and disk-related alerting.

The workload uses the same validation and safety limits as the HTTP API.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		req := app.DiskStressRequest{
			Size:     stressDiskOpts.size,
			Duration: stressDiskOpts.duration,
		}
		result, err := app.RunDiskStress(cmd.Context(), req)
		if err != nil {
			return err
		}

		cmd.Printf("done: size=%s bytes=%d duration=%s\n", result.Size, result.Bytes, result.Duration)

		return nil
	},
}

func init() {
	stressCmd.AddCommand(stressDiskCmd)

	stressDiskCmd.Flags().StringVar(&stressDiskOpts.size, "size", "", "disk size to write")
	stressDiskCmd.Flags().StringVar(&stressDiskOpts.duration, "duration", "", "stress duration")

	if err := stressDiskCmd.MarkFlagRequired("duration"); err != nil {
		panic(err)
	}
	if err := stressDiskCmd.MarkFlagRequired("size"); err != nil {
		panic(err)
	}
}
