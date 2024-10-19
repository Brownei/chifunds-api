package main

import (
	"database/sql"
	"log"

	"github.com/brownei/chifunds-api/cmd/api"
	"github.com/brownei/chifunds-api/db"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	newDb, err := db.NewPostgresStorage()
	if err != nil {
		log.Fatalf("Database connection unsuccessful: %s", err)
	}

	initializeDb(newDb)
	db.AddMigrations(newDb)

	server := api.NewServer(":8000", newDb)
	if err := server.Run(); err != nil {
		log.Printf("Cannot start up server: %s", err)
	}
}

func initializeDb(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("Couldn not connect to database: %s", err)
	}

	log.Printf("Database connected")
}
