package auth

import (
	"context"
	"net/http"

	db "github.com/brownei/chifunds-api/db"
	"github.com/brownei/chifunds-api/service"
	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	store types.AuthStore
}

func NewAuthHandler(store types.AuthStore) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

func (a *AuthHandler) AllAuthRoutes(r chi.Router) {

	r.Get("/{provider}", a.GoogleAuthLoginAndRegister)
	r.Get("/{provider}/callback", a.ProviderAuthCallbackFunction)
}

func (a *AuthHandler) GoogleAuthLoginAndRegister(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func (a *AuthHandler) ProviderAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	var u *types.User
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	ctx := r.Context()

	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		db, _ := db.NewPostgresStorage()
		store := service.NewStore(db)

		existingUSer, err := store.GetUsersByEmail(ctx, gothUser.Email)
		if err != nil {
			hashedPassword, err := HashPassword(gothUser.Email)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, err)
			}

			u, err = store.CreateNewUser(ctx, types.RegisterUserPayload{
				Email:          gothUser.Email,
				FirstName:      gothUser.FirstName,
				LastName:       gothUser.LastName,
				ProfilePicture: gothUser.AvatarURL,
				Password:       hashedPassword,
				EmailVerified:  true,
			})

			utils.WriteJSON(w, http.StatusOK, u)
		} else if existingUSer != nil {
			utils.WriteJSON(w, http.StatusOK, existingUSer)
		}

	} else {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

}
