package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	Validator = validator.New()
)

func VerifyPassword(encryptedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)); err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

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

func RandomAccountNumber() string {
	rand := 1000000 + rand.Intn(9999999-1000000)
	return fmt.Sprintf("512%d", rand)
}

func RandomPhoneNumber() string {
	ran := 100000000 + rand.Intn(999999999-100000000)
	return fmt.Sprintf("+23480%s", ran)
}

func JwtToken(email string, ctx context.Context) string {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	expiryTime := time.Hour * 168
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,                     // Subject (user identifier)
		"iss": "chifunds",                // Issuer
		"exp": expiryTime.Milliseconds(), // Expiration time
		"iat": time.Now().Unix(),         // Issued at
	})

	token, _ := claims.SignedString(secretKey)
	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", token)
	return token
}

func VerifyToken(token string) (string, error) {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	verifiedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("Error in the verified token: %s", err.Error())
	}

	// Check if the token is valid
	if !verifiedToken.Valid {
		return "", fmt.Errorf("Not Valid!", err.Error())
	}

	email, _ := verifiedToken.Claims.GetSubject()
	// Return the verified token
	log.Printf("VerifiedToken: %v\n", email)

	return email, nil
}
