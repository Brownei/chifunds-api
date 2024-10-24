package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/markbates/goth/gothic"
)

func (a *application) AllAuthRoutes(r chi.Router) {
	r.Post("/signin", a.Login)
	r.Post("/signup", a.CreateAUser)
	r.Get("/{provider}", a.GoogleAuthLoginAndRegister)
	r.Get("/{provider}/callback", a.ProviderAuthCallbackFunction)
}

func (a *application) GoogleAuthLoginAndRegister(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func (a *application) CreateAUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := utils.ParseJSON(r, &payload); err != nil {
		log.Printf("Error: %s", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//Validate the payload
	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		fmt.Printf("Error: %s", errors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", errors))
		return
	}

	//Check if user exists first
	existingUser, err := a.store.Users.GetUsersByEmail(ctx, payload.Email, false)
	log.Printf("existingUser: %v\n", existingUser)
	if existingUser != nil {
		utils.WriteError(w, http.StatusFound, fmt.Errorf("User already exists"))
		return
	}

	_, err = a.store.Users.CreateNewUser(ctx, types.RegisterUserPayload{
		Email:          payload.Email,
		FirstName:      payload.FirstName,
		Password:       payload.Password,
		LastName:       payload.LastName,
		EmailVerified:  false,
		ProfilePicture: payload.ProfilePicture,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	token := utils.JwtToken(payload.Email, ctx)
	log.Printf("Token: %s", token)
	utils.WriteJSON(w, http.StatusCreated, token)
}

func (a *application) ProviderAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	ctx := r.Context()

	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		existingUSer, err := a.store.Users.GetUsersByEmail(ctx, gothUser.Email, false)
		if err != nil {
			hashedPassword, err := HashPassword(gothUser.Email)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, err)
			}

			_, err = a.store.Users.CreateNewUser(ctx, types.RegisterUserPayload{
				Email:          gothUser.Email,
				FirstName:      gothUser.FirstName,
				LastName:       gothUser.LastName,
				ProfilePicture: gothUser.AvatarURL,
				Password:       hashedPassword,
				EmailVerified:  true,
			})

			token := utils.JwtToken(gothUser.Email, ctx)
			log.Printf("Token: %s", token)
			utils.WriteJSON(w, http.StatusOK, token)
		} else if existingUSer != nil {
			utils.WriteJSON(w, http.StatusOK, existingUSer)
		}

	} else {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
}

func (a *application) Login(w http.ResponseWriter, r *http.Request) {
	var loginPayload types.LoginPayload
	ctx := r.Context()
	if err := utils.ParseJSON(r, &loginPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	existingUser, err := a.store.Users.GetUsersByEmail(ctx, loginPayload.Email, true)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	} else if existingUser == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("No user like this exists"))
	}

	token, err := a.store.Auth.Login(ctx, existingUser, loginPayload)

	utils.WriteJSON(w, http.StatusAccepted, token)
}
