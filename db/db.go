package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	name     = os.Getenv("DB_NAME")
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
)

func NewPostgresStorage() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	db, err := sql.Open("postgres", connStr)

	return db, err
}

func InitializeDb(db *sql.DB, logger *zap.SugaredLogger) {
	err := db.Ping()
	if err != nil {
		logger.Fatalf("Couldn not connect to database: %s", err)
	}

	logger.Info("Database connected")
}
