package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const name = "greenlight"

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func newRootCommand(config *config) *cobra.Command {
	rootCmd := cobra.Command{
		Use:   name,
		Short: name,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	rootCmd.PersistentFlags().StringVarP(&config.env, "env", "e", "dev", "Env file path")
	rootCmd.PersistentFlags().IntVarP(&config.port, "port", "p", 8080, "Server port")
	return &rootCmd
}

func main() {
	var cfg config
	rootCmd := newRootCommand(&cfg)
	rootCmd.Execute()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		config: cfg,
		logger: logger,
	}

	router := app.routes()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
