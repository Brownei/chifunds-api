package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var (
	Validator = validator.New()
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		log.Printf("Missing body data")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func RandomBVN() int {
	return 10000000000 + rand.Intn(99999999999-10000000000)
}

func RandomPhoneNumber() string {
	ran := 10000000 + rand.Intn(99999999-10000000)
	return fmt.Sprintf("+23480%s", ran)
}
