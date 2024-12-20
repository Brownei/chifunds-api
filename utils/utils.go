package utils

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/brownei/chifunds-api/types"
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

func DecryptAndParseJson(r *http.Request, decryptFunc func(string) ([]byte, error)) ([]byte, error) {
	var payload types.DataPayload
	if r.Body == nil {
		log.Printf("Missing body data")
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}

	decryptedData, err := decryptFunc(payload.Data)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func EncryptAndWriteJson(w http.ResponseWriter, status int, byteTrans []byte, encryptFunc func([]byte) (string, string, error)) {
	encrypted, key, err := encryptFunc(byteTrans)
	if err != nil {
		WriteError(w, http.StatusBadGateway, err)
	}

	payload := types.EncryptedDataPayload{
		Encrypted: encrypted,
		AesKey:    key,
	}

	WriteJSON(w, status, payload)
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

func JwtToken(email string, ctx context.Context) string {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	expiryTime := time.Now().Add(7 * 24 * time.Hour).Unix() // 7 days in seconds
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,      // Subject (user identifier)
		"iss": "chifunds", // Issuer
		"exp": expiryTime,
		"iat": time.Now().Unix(), // Issued at
	})

	token, _ := claims.SignedString(secretKey)
	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", token)
	return token
}

func VerifyToken(token string) (string, error) {
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	verifiedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Error: %v", ok)
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("Error in the verified token: %s\n%v", err.Error(), verifiedToken)
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

func StructToBytes(data any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
