package store

import (
	"context"
	"database/sql"
	"log"

	"github.com/brownei/chifunds-api/types"
	"go.uber.org/zap"
)

type Store struct {
	Users interface {
		GetChifundsUser(email string, forLogin bool) (*types.User, error)
		GetUsersByEmail(ctx context.Context, email string, forLogn bool) (*types.User, error)
		GetAllUsers() ([]types.User, error)
		CreateNewUser(ctx context.Context, payload types.RegisterUserPayload) (*types.User, error)
		CreateChiFundsAdminUser(payload types.RegisterUserPayload) error
		GetBalance(context.Context, string) (*types.Balance, error)
	}

	Auth interface {
		Login(ctx context.Context, existingUser *types.User, payload types.LoginPayload) (string, error)
	}

	Transactions interface {
		BorrowMoney(ctx context.Context, logger *zap.SugaredLogger, lendedMoney int32, userId int8) error
		TransferMoney(context.Context, *zap.SugaredLogger, types.User, int32, string) error
		GetReceivedTransactions(context.Context, string) ([]types.ReceivedTransactions, error)
		GetSentTransactions(context.Context, string) ([]types.SentTransactions, error)
		GetBorrowedTransactions(context.Context) ([]types.BorrowedTransactions, error)
	}
}

var (
	users []types.User
)

func NewStore(db *sql.DB) Store {
	return Store{
		Users:        &UserStore{db},
		Auth:         &AuthStore{db},
		Transactions: &TransactionStore{db},
	}
}

func (s *Store) CreateChiFundsUser() error {
	payload := types.RegisterUserPayload{
		Email:          "chifundsadmin@gmail.com",
		FirstName:      "ChiFunds",
		LastName:       "Funding",
		ProfilePicture: "",
		Password:       "sfhkbhagassvnldfhdgklhdhguytigndnb",
		EmailVerified:  true,
	}

	_, err := s.Users.GetChifundsUser(payload.Email, false)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := s.Users.CreateChiFundsAdminUser(payload); err != nil {
				if err == sql.ErrNoRows {
					return nil
				}

				return err
			}

			log.Printf("Admin user created!")
			return nil

		} else {
			log.Printf("Admin user already available")
			return nil
		}
	}

	return nil
}
