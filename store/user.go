package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) GetChifundsUser(email string, forLogin bool) (*types.User, error) {
	query := `SELECT id FROM "user" WHERE email = $1`
	u := &types.User{}

	err := s.db.QueryRow(query, email).Scan(
		&u.ID,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserStore) GetUsersByEmail(ctx context.Context, email string, forLogin bool) (*types.User, error) {
	var query string
	//wg := sync.WaitGroup{}
	if forLogin == true {
		query = `SELECT u.id, u.email, u.first_name, u.last_name, u.profile_picture, u.email_verified, u.password, a.account_number, a.money FROM "user" AS u JOIN "account" AS a ON u.id = a.user_id WHERE u.email = $1`
	} else {
		query = `SELECT u.id, u.email, u.first_name, u.last_name, u.profile_picture, u.email_verified, a.account_number, a.money FROM "user" AS u JOIN "account" AS a ON u.id = a.user_id WHERE u.email = $1`
	}

	u := &types.User{}

	if forLogin == true {
		err := s.db.QueryRowContext(ctx, query, email).Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.ProfilePicture,
			&u.EmailVerified,
			&u.Password,
			&u.AccountNumber,
			&u.Balance,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("No user like this!")
			}

			return nil, err
		}

		return u, nil

	} else {
		err := s.db.QueryRowContext(ctx, query, email).Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.ProfilePicture,
			&u.EmailVerified,
			&u.AccountNumber,
			&u.Balance,
		)

		if err != nil {
			return nil, fmt.Errorf("No user like this!")
		}

		return u, nil
	}
}

func (s *UserStore) GetAllUsers() ([]types.User, error) {
	var u []types.User
	query := "SELECT id, email, first_name, last_name, profile_picture, email_verified FROM \"user\" "

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user types.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.ProfilePicture,
			&user.EmailVerified,
		)
		if err != nil {
			return nil, err
		}

		u = append(u, user)
	}

	return u, nil
}

func (s *UserStore) CreateChiFundsAdminUser(payload types.RegisterUserPayload) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		log.Printf("Couldn't hash a password: %s", err)
		return err
	}

	creatingNewUserQuery := `INSERT INTO "user" (email, first_name, last_name, profile_picture, password, email_verified) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = s.db.Query(creatingNewUserQuery, payload.Email, payload.FirstName, payload.LastName, payload.ProfilePicture, hashPassword, payload.EmailVerified)

	return err
}

func (s *UserStore) CreateNewUser(ctx context.Context, payload types.RegisterUserPayload) (*types.User, error) {
	user := &types.User{}
	accountNumber := utils.RandomAccountNumber()
	log.Printf("account: %s", accountNumber)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		log.Printf("Couldn't hash a password: %s", err)
		return nil, err
	}

	existingAccountNumberQuery := `SELECT id FROM "account" WHERE account_number = $1`

	err = s.db.QueryRowContext(ctx, existingAccountNumberQuery, accountNumber).Scan()
	if err == nil {
		return nil, fmt.Errorf("Account number unavaialble")
	}

	creatingNewUserQuery := `INSERT INTO "user" (email, first_name, last_name, profile_picture, password, email_verified) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, email, first_name, last_name, profile_picture, email_verified`

	err = s.db.QueryRowContext(ctx, creatingNewUserQuery, payload.Email, payload.FirstName, payload.LastName, payload.ProfilePicture, hashPassword, payload.EmailVerified).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.ProfilePicture,
		&user.EmailVerified,
	)

	creatingNewAccountQuery := `INSERT INTO "account" (account_number, user_id) VALUES ($1, $2) RETURNING account_number`
	err = s.db.QueryRowContext(ctx, creatingNewAccountQuery, accountNumber, user.ID).Scan(
		&user.AccountNumber,
	)
	//_, err = scanRowsToReturnUser(rows)
	return user, err
}

func (s *UserStore) GetBalance(ctx context.Context, email string) (*types.Balance, error) {
	balance := &types.Balance{}
	query := `SELECT a.money FROM "account" AS a JOIN "user" AS u ON u.id = a.user_id WHERE u.email = $1`

	if err := s.db.QueryRowContext(ctx, query, email).Scan(
		&balance.Amount,
	); err != nil {
		return nil, err
	}

	return balance, nil
}
