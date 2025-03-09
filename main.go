package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/wildan3105/converto/cmd"
	config "github.com/wildan3105/converto/configs"
)

func main() {
	config.LoadConfig()

	rootCmd := &cobra.Command{
		Use:   "engine",
		Short: "converto command line interface",
	}

	rootCmd.AddCommand(cmd.RestCmd)
	rootCmd.AddCommand(cmd.WorkerCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
		os.Exit(1)
	}
}
