package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.UserStore
}

var (
	payload types.RegisterUserPayload
)

func NewUserHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) AllUsersRoutes(r chi.Router) {
	r.Get("/", h.GetAllUsers)
	r.Post("/", h.CreateAUser)
	r.Get("/{email}", h.GetUsersByEmail)
}

func (h *Handler) GetUsersByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	ctx := r.Context()

	u, err := h.store.GetUsersByEmail(ctx, email)
	if err != nil {
		fmt.Errorf("Error getting user: %v\n", err)
		utils.WriteError(w, http.StatusBadRequest, err)
	} else if u == nil {
		utils.WriteJSON(w, http.StatusNotFound, fmt.Errorf("No user like this!"))
	}

	fmt.Printf("User: %v\n", u)
	utils.WriteJSON(w, http.StatusOK, u)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	u, err := h.store.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error in getting all users: %s", err))
	}

	utils.WriteJSON(w, http.StatusOK, u)
	return
}

func (h *Handler) CreateAUser(w http.ResponseWriter, r *http.Request) {
	//Check if there is a body
	ctx := r.Context()
	//var res map[string]interface{}
	//body := r.Body

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
	existingUser, err := h.store.GetUsersByEmail(ctx, payload.Email)
	if existingUser != nil {
		utils.WriteError(w, http.StatusFound, fmt.Errorf("User already exists"))
		return
	}

	user, err := h.store.CreateNewUser(ctx, types.RegisterUserPayload{
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

	//log.Printf("Response: ", responseReader)
	utils.WriteJSON(w, http.StatusCreated, user)
}

func (a *Handler) GoogleAuthLoginAndRegister(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Success!"))
}
