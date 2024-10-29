package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/brownei/chifunds-api/utils"
)

func (a *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		a.logger.Info(authHeader)
		if authHeader == "" {
			a.logger.Errorf("Unauthorized permission: %s", fmt.Errorf("No token"))
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("No token"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.logger.Errorf("Unauthorized permission: %s", fmt.Errorf("No Bearer token"))
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Not the Bearer Token"))
			return
		}

		token := parts[1]
		userEmail, err := utils.VerifyToken(token)
		if err != nil {
			//a.logger.Errorf("Unauthorized permission: %s", err.Error())
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, "user", userEmail)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *application) PublicKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretToken := os.Getenv("SECRET_KEY")
		authHeader := r.Header.Get("Authorization")
		a.logger.Info(authHeader)
		if authHeader == "" {
			a.logger.Errorf("Unauthorized permission: %s", fmt.Errorf("No token"))
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("No token"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.logger.Errorf("Unauthorized permission: %s", fmt.Errorf("No Bearer token"))
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Not the Bearer Token"))
			return
		}

		token := parts[1]
		a.logger.Info(secretToken, token)
		if token != secretToken {
			//a.logger.Errorf("Unauthorized permission: %s", err.Error())
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized to view call this method"))
			return
		}

		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
