package main

import (
	"log"

	"github.com/brownei/chifunds-api/cmd/api"
	"github.com/brownei/chifunds-api/db"
	"github.com/brownei/chifunds-api/store"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	logger := zapLogger.Sugar()

	defer logger.Sync()

	newDb, err := db.NewPostgresStorage()
	store := store.NewStore(newDb)
	if err != nil {
		logger.Error("Database connection unsuccessful: %s", zap.Field{
			Interface: err,
		})
	}

	db.InitializeDb(newDb, logger)
	db.AddMigrations(newDb, logger)
	if err := api.GenerateRSAKeys(); err != nil {
		log.Printf("Error generating keys: %v", err)
	}

	server := api.NewServer(":8000", logger, newDb, store)
	if err := server.Run(); err != nil {
		log.Printf("Cannot start up server: %s", err)
	}
}
