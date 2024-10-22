package store

import (
	"context"
	"database/sql"

	"github.com/brownei/chifunds-api/types"
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
}

var (
	users []types.User
)

func NewStore(db *sql.DB) Store {
	return Store{
		Users: &UserStore{db},
		Auth:  &AuthStore{db},
	}
}
