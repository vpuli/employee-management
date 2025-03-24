package main

import (
	"employees/api"
	"employees/internal/db/postgres"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger() (*zap.Logger, error) {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Create log file with current date
	currentTime := time.Now()
	logFile := filepath.Join(logsDir, fmt.Sprintf("%s.log", currentTime.Format("2006-01-02")))

	// Open the log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		zap.InfoLevel,
	)

	// Create logger
	logger := zap.New(core)
	return logger, nil
}

func main() {
	fmt.Println("Server starting...")

	// Setup logger
	logger, err := setupLogger()
	if err != nil {
		fmt.Printf("Failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		logger.Fatal("Failed to get current user", zap.Error(err))
	}

	// Initialize database
	dbConfig := postgres.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     currentUser.Username,
		Password: "",
		DBName:   "employees",
		SSLMode:  "disable",
	}

	// Update the environment to use PostgreSQL 15
	os.Setenv("PGHOST", dbConfig.Host)
	os.Setenv("PGPORT", dbConfig.Port)
	os.Setenv("PGUSER", dbConfig.User)
	os.Setenv("PGDATABASE", dbConfig.DBName)

	db, err := postgres.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func() {
		logger.Sync()
		db.Close()
	}()

	s := api.NewServer(":8080", logger, db)

	fmt.Println("Server listening on :8080")
	logger.Fatal("Server error", zap.Error(s.Start()))
}
