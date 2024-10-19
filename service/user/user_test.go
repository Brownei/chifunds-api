package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brownei/chifunds-api/types"
	"github.com/go-chi/chi/v5"
)

type mockUserStore struct {
}

func TestUserServiceHandler(t *testing.T) {
	userStore := &mockUserStore{}
	userHandler := NewHandler(userStore)

	t.Run("should create a new user successfully", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			Email:          "brownson@gmail.com",
			FirstName:      "Email",
			LastName:       "Testing",
			ProfilePicture: "",
			Password:       "12345",
			EmailVerified:  false,
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/users", userHandler.CreateAUser)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code should be %d, but got %d", http.StatusOK, rr.Code)
		}
	})
}

func (m *mockUserStore) GetUsersByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("User not found")
}
func (m *mockUserStore) CreateNewUser(payload types.RegisterUserPayload) error {
	return nil
}

func (m *mockUserStore) GetAllUsers() (*[]types.User, error) {
	return nil, fmt.Errorf("No users")
}
