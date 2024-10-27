package db

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
)

func AddMigrations(db *sql.DB, log *zap.SugaredLogger) {
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
					`CREATE TABLE IF NOT EXISTS "card" (id SERIAL PRIMARY KEY, serial_no VARCHAR(11) NOT NULL UNIQUE, cvc INT, expiry_date TIMESTAMP, account_id INT NOT NULL)`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "card"`,
				},
			},

			{
				Id: "3",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "account" (id SERIAL PRIMARY KEY, account_number VARCHAR(10), money INT DEFAULT 0, user_id INT NOT NULL, card_id INT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "account"`,
				},
			},

			{
				Id: "4",
				Up: []string{
					`ALTER TABLE "card" ADD CONSTRAINT "Card_accountId_fkey" FOREIGN KEY ("account_id") REFERENCES "account"("id")`,
				},
				Down: []string{
					`ALTER TABLE "card" DROP CONSTRAINT IF EXISTS "Card_accountId_fkey"`,
				},
			},

			{
				Id: "5",
				Up: []string{
					`ALTER TABLE "account" ADD CONSTRAINT "Account_userId_fkey" FOREIGN KEY ("user_id") REFERENCES "user"("id")`,
				},
				Down: []string{
					`ALTER TABLE "account" DROP CONSTRAINT IF EXISTS "Account_userId_fkey"`,
				},
			},

			{
				Id: "6",
				Up: []string{
					`ALTER TABLE "account" ADD CONSTRAINT "Account_cardId_fkey" FOREIGN KEY ("card_id") REFERENCES "card"("id")`,
				},
				Down: []string{
					`ALTER TABLE "account" DROP CONSTRAINT IF EXISTS "Account_cardId_fkey"`,
				},
			},

			{
				Id: "7",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS "transactions" (id SERIAL PRIMARY KEY, receiver_id INT, sender_id INT, amount_sent INT, sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`,
				},
				Down: []string{
					`DROP TABLE IF EXISTS "transactions"`,
				},
			},

			{
				Id: "8",
				Up: []string{
					`ALTER TABLE "transactions" ADD CONSTRAINT "Transactions_receiverId_fkey" FOREIGN KEY ("receiver_id") REFERENCES "user"("id")`,
				},
				Down: []string{
					`ALTER TABLE "transactions" DROP CONSTRAINT IF EXISTS "Transactions_receiverId_fkey"`,
				},
			},

			{
				Id: "9",
				Up: []string{
					`ALTER TABLE "transactions" ADD CONSTRAINT "Transactions_senderId_fkey" FOREIGN KEY ("sender_id") REFERENCES "user"("id")`,
				},
				Down: []string{
					`ALTER TABLE "transactions" DROP CONSTRAINT IF EXISTS "Transactions_senderId_fkey"`,
				},
			},
		},
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Couldn't apply the migrations: %s", err)
	}

	log.Infof("Applied %d migrations!", n)
}
