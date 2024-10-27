package api

import (
	"context"
	"fmt"
	"net/http"
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
