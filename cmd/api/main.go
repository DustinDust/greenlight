package main

import (
	"fmt"
	"greenlight/internal/data"
	"greenlight/internal/jsonlog"
	"log"
	"net/http"
	"os"
	"time"

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

	router := app.routes()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      router,
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.PrintInfo(fmt.Sprintf("starting %s server on %s", cfg.env, srv.Addr), map[string]string{"env": cfg.env, "addr": srv.Addr})
	err = srv.ListenAndServe()
	logger.PrintFatal(err.Error(), nil)
}
