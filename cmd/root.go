package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "rwth-calendar",
	}

	rootCmd.AddCommand(NewGenerateCmd())
	rootCmd.AddCommand(NewServeCmd())

	return rootCmd
}
