package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alkindi42/probelet/internal/engine"
	apphttp "github.com/Alkindi42/probelet/internal/http"
	"github.com/spf13/cobra"
)

// serveOptions holds the flags for the serve command.
type serveOptions struct {
	initialReadinessDelay time.Duration
}

var serveOpts serveOptions

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Probelet HTTP API",
	Long:  `Launch the HTTP server to expose endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)

		addr := ":8000"
		readinessStore := engine.NewReadinessStore()
		readinessCtl := engine.NewReadinessController(readinessStore)

		runCtx, runCancel := context.WithCancel(context.Background())
		defer runCancel()

		if serveOpts.initialReadinessDelay > 0 {
			readinessCtl.StartInitialDelay(runCtx, serveOpts.initialReadinessDelay)
		}

		router := apphttp.NewRouter(readinessCtl)
		server := &http.Server{
			Addr:    addr,
			Handler: router,
		}

		go func() {
			slog.Info("starting HTTP server", "addr", addr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("http server error", "err", err)
				os.Exit(1)
			}
		}()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		slog.Info("shutdown signal received")

		runCancel()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("forced shutdown", "err", err)
			os.Exit(1)
		}

		slog.Info("server stopped gracefully")
	},
}

func init() {
	serveCmd.Flags().DurationVar(
		&serveOpts.initialReadinessDelay,
		"initial-readiness-delay",
		0,
		"Time to keep the server in non-ready state before becoming ready (e.g. 10s, 5m)",
	)
	rootCmd.AddCommand(serveCmd)
}
