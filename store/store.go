package store

import (
	"context"
	"database/sql"

	"github.com/brownei/chifunds-api/types"
	"go.uber.org/zap"
)

type Store struct {
	Users interface {
		GetUsersByEmail(ctx context.Context, email string, forLogn bool) (*types.User, error)
		GetAllUsers() ([]types.User, error)
		CreateNewUser(ctx context.Context, payload types.RegisterUserPayload) (*types.User, error)
	}

	Auth interface {
		Login(ctx context.Context, existingUser *types.User, payload types.LoginPayload) (string, error)
	}

	Transactions interface {
		BorrowMoney(ctx context.Context, lendedMoney int32, userId int8) error
		TransferMoney(context.Context, *zap.SugaredLogger, types.User, int32, string) error
		GetReceivedTransactions(context.Context, string) (*types.ReceivedTransactions, error)
		GetSentTransactions(context.Context, string) (*types.SentTransactions, error)
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
