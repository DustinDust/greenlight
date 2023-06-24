package main

import (
	"os"

	"github.com/spf13/cobra"
)

func newRootCommand(config *config) *cobra.Command {
	rootCmd := cobra.Command{
		Use:   name,
		Short: name,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	rootCmd.PersistentFlags().StringVarP(&config.env, "env", "e", "dev", "Environment (dev|prod)")
	rootCmd.PersistentFlags().IntVarP(&config.port, "port", "p", 8080, "Server port")
	rootCmd.PersistentFlags().StringVarP(&config.db.dsn, "db-dsn", "d", os.Getenv("GREENLIGHT_DB_CONNECTION"), "PostgresSQL dsn")
	rootCmd.PersistentFlags().IntVar(&config.db.maxIdleConns, "db-max-idle-conns", 25, "Maximum number of idle connections")
	rootCmd.PersistentFlags().IntVar(&config.db.maxOpenConns, "db-max-open-conns", 25, "Maximum number of open connections")
	rootCmd.PersistentFlags().StringVar(&config.db.maxIdleTime, "db-max-idle-time", "15m", "Maximum duration of an idle connection")

	// Exit the program after help template is printed out
	defaultHelpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(c *cobra.Command, s []string) {
		defaultHelpFunc(c, s)
		os.Exit(0)
	})
	return &rootCmd
}
