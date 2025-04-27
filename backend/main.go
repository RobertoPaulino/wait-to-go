package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"wait-to-go/auth"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func loadConfig() (*Config, error) {
	config := &Config{
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "sicreto"),
		DBName:     getEnvOrDefault("DB_NAME", "gopgtest"),
		DBSSLMode:  getEnvOrDefault("DB_SSL_MODE", "disable"),
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize database
	if err = createEntryTable(db); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Initialize queue and history
	entryQueue := []Entry{}
	historySlice := []Entry{}

	// Load waiting entries from database
	waitingEntries, err := getWaitingEntry(db)
	if err != nil {
		log.Printf("Warning: Failed to load waiting entries: %v", err)
	} else {
		entryQueue = waitingEntries
	}

	app := App{
		db:      db,
		queue:   &entryQueue,
		history: &historySlice,
	}

	// Setup routes with CORS and authentication middleware
	mux := http.NewServeMux()

	// Public endpoint
	mux.HandleFunc("/join", enableCors(app.handleJoin))

	// Customer endpoints (require JWT)
	mux.HandleFunc("/status/", enableCors(auth.AuthMiddleware(app.handleStatus)))

	// Admin endpoints (require API key)
	mux.HandleFunc("/queue", enableCors(auth.AdminAuthMiddleware(app.handleQueue)))
	mux.HandleFunc("/next", enableCors(auth.AdminAuthMiddleware(app.handleNext)))
	mux.HandleFunc("/serve", enableCors(auth.AdminAuthMiddleware(app.handleServe)))
	mux.HandleFunc("/clear", enableCors(auth.AdminAuthMiddleware(app.handleClear)))

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
