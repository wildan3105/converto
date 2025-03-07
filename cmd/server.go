package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	config "github.com/wildan3105/converto/configs"
	"github.com/wildan3105/converto/pkg/api"
	"github.com/wildan3105/converto/pkg/logger"
)

// RestCmd is the command to run the REST API server
var RestCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the REST API server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.GetInstance()

		logger.Info("Starting REST API server...")

		app := api.Setup()

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			if err := app.Listen(":" + config.AppConfig.Port); err != nil {
				logger.Error("Server failed: %v", err)
			}
		}()

		sig := <-done
		logger.Info("Signal received: %s. Shutting down server gracefully...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			logger.Error("Server shutdown failed: %v", err)
		}

		logger.Info("Server exited gracefully")
	},
}
