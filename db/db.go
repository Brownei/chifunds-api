package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

var (
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	name     = os.Getenv("DB_NAME")
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
)

func NewPostgresStorage() (*sql.DB, error) {
	log.Printf(user)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	db, err := sql.Open("postgres", connStr)

	return db, err
}

func AddMigrations(db *sql.DB) {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{

			{
				Id: "1",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "user" (id SERIAL PRIMARY KEY, email VARCHAR(100) NOT NULL UNIQUE, first_name VARCHAR(100) NOT NULL, last_name VARCHAR(100) NOT NULL, profile_picture VARCHAR(255), password VARCHAR(100) NOT NULL, email_verified BOOLEAN DEFAULT FALSE, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
				},
				Down: []string{`"DROP TABLE IF EXISTS "user"`},
			},

			{
				Id: "2",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "account_number" (id SERIAL PRIMARY KEY, subaccount_id VARCHAR(100) NOT NULL, subaccount_number VARCHAR(100) NOT NULL, wallet_id VARCHAR(100) NOT NULL, wallet_number VARCHAR(100) NOT NULL, user_id INT NOT NULL);`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "account_number"`,
				},
			},
			{
				Id: "3",
				Up: []string{
					`ALTER TABLE "account_number" ADD CONSTRAINT "AccountNumber_userId_fkey" FOREIGN KEY ("user_id") REFERENCES "user"("id")`,
				},
				Down: []string{
					`ALTER TABLE "account_number" DROP CONSTRAINT IF EXISTS "AccountNumber_userId_fkey";`,
				}},
		},
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Couldn't apply the migrations: %s", err)
	}

	log.Printf("Applied %d migrations!\n", n)
}
