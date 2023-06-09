package main

import "github.com/spf13/cobra"

func NewRoot(config *config) *cobra.Command {
	rootCmd := cobra.Command{
		Use:   name,
		Short: name,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	rootCmd.PersistentFlags().StringVarP(&config.env, "env", "e", ".env", "Env file path")
	rootCmd.PersistentFlags().IntVarP(&config.port, "port", "p", 8080, "Server port")
	return &rootCmd
}
