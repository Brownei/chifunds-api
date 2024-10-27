package api

import (
	"fmt"
	"net/http"

	"github.com/brownei/chifunds-api/utils"
	"github.com/go-chi/chi/v5"
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
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	} else if u == nil {
		utils.WriteJSON(w, http.StatusNotFound, fmt.Errorf("No user like this!"))
		return
	}

	fmt.Printf("User: %v\n", u)
	utils.WriteJSON(w, http.StatusOK, u)
	return
}

func (a *application) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	u, err := a.store.Users.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error in getting all users: %s", err))
	}

	utils.WriteJSON(w, http.StatusOK, u)
	return
}
