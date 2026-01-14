package cmd

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apphttp "github.com/Alkindi42/probelet/internal/http"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Probelet HTTP API",
	Long:  `Launch the HTTP server to expose endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)

		addr := ":8000"
		router := apphttp.NewRouter()

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
	rootCmd.AddCommand(serveCmd)
}
