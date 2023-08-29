package main

import (
	"greenlight/internal/data"
	"greenlight/internal/jsonlog"
	"os"

	_ "github.com/lib/pq"
)

const name = "greenlight"

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	rootCmd := newRootCommand(&cfg)
	rootCmd.Execute()
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal("error opening connection to database: %v", map[string]string{"err": err.Error()})
	}
	logger.PrintInfo("database connection pool established", nil)
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err.Error(), nil)
	}
}
