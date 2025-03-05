package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grpc",
	Short: "A gRPC server and client CLI",
}

func main() {
	// Add commands to rootCmd
	rootCmd.AddCommand(ServerCmd)
	rootCmd.AddCommand(ClientCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
