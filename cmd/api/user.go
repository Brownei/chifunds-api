package api

import (
	"fmt"
	"net/http"

	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/utils"
	"github.com/go-chi/chi/v5"
)

var (
	payload types.RegisterUserPayload
)

func (a *application) AllUsersRoutes(r chi.Router) {
	r.Get("/", a.GetAllUsers)
	r.Get("/{email}", a.GetUsersByEmail)
}

func (a *application) GetUsersByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	ctx := r.Context()

	u, err := a.store.Users.GetUsersByEmail(ctx, email, false)
	if err != nil {
		fmt.Errorf("Error getting user: %v\n", err)
		utils.WriteError(w, http.StatusBadRequest, err)
	} else if u == nil {
		utils.WriteJSON(w, http.StatusNotFound, fmt.Errorf("No user like this!"))
	}

	fmt.Printf("User: %v\n", u)
	utils.WriteJSON(w, http.StatusOK, u)
}

func (a *application) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	u, err := a.store.Users.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error in getting all users: %s", err))
	}

	utils.WriteJSON(w, http.StatusOK, u)
	return
}
