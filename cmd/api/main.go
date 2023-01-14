package main

import (
	"context"      // New import
	"database/sql" // New import
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// Import the pq driver so that it can register itself with the database/sql
	// package. Note that we alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	// _ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"greenlight.alexedwards.net/internal/data"
	"greenlight.alexedwards.net/internal/jsonlog"
)
const version = "1.0.0"
// Add maxOpenConns, maxIdleConns and maxIdleTime fields to hold the configuration
// settings for the connection pool.
type config struct {
	port int
	env  string
	db   struct {
			dsn          string
			maxOpenConns int
			maxIdleConns int
			maxIdleTime  string
	}
}
// Add a models field to hold our new Models struct.
// Change the logger field to have the type *jsonlog.Logger, instead of
// *log.Logger.
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()
	// Initialize a new jsonlog.Logger which writes any messages *at or above* the INFO
	// severity level to the standard out stream.
	// Initialize the custom logger.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
		
	db, err := openDB(cfg)
	if err != nil {
			// Use the PrintFatal() method to write a log entry containing the error at the
			// FATAL level and exit. We have no additional properties to include in the log
			// entry, so we pass nil as the second parameter.
			logger.PrintFatal(err, nil)
	}
	defer db.Close()
	// Likewise use the PrintInfo() method to write a message at the INFO level.
	logger.PrintInfo("database connection pool established", nil)
	app := &application{
			config: cfg,
			logger: logger,
			models: data.NewModels(db),
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		// Create a new Go log.Logger instance with the log.New() function, passing in
		// our custom Logger as the first parameter. The "" and 0 indicate that the
		// log.Logger instance should not use a prefix or any flags.
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
}
	// Again, we use the PrintInfo() method to write a "starting server" message at the
	// INFO level. But this time we pass a map containing additional properties (the
	// operating environment and server address) as the final parameter.
	logger.PrintInfo("starting server", map[string]string{
			"addr": srv.Addr,
			"env":  cfg.env,
	})
	err = srv.ListenAndServe()
	// Use the PrintFatal() method to log the error and exit.
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
			return nil, err
	}
	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
			return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
			return nil, err
	}
	return db, nil
}