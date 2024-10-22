package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/utils"
)

type AuthStore struct {
	store *sql.DB
}

func (s *AuthStore) Login(ctx context.Context, existingUser *types.User, payload types.LoginPayload) (string, error) {
	//Get the user from the database
	if err := utils.VerifyPassword(existingUser.Password, payload.Password); err != nil {
		return "Not verified", fmt.Errorf(err.Error())
	}

	token := utils.JwtToken(payload.Email, ctx)

	return token, nil
}
