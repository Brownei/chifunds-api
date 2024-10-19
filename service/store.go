package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	db *sql.DB
}

var (
	users []types.User
)

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUsersByEmail(ctx context.Context, email string) (*types.User, error) {
	query := "SELECT id, email, first_name, last_name, profile_picture, email_verified FROM \"user\" WHERE email = $1"
	u := &types.User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.ProfilePicture,
		&u.EmailVerified,
	)

	if err != nil {
		return nil, fmt.Errorf("No user like this!")
	}

	return u, nil
}

func (s *Store) GetAllUsers() ([]types.User, error) {
	var u []types.User
	query := "SELECT id, email, first_name, last_name, profile_picture, email_verified FROM \"user\" "

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user types.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.ProfilePicture,
			&user.EmailVerified,
		)
		if err != nil {
			return nil, err
		}

		u = append(u, user)
	}

	return u, nil
}

func (s *Store) CreateNewUser(ctx context.Context, payload types.RegisterUserPayload) (*types.User, error) {
	user := &types.User{}
	walletChan := make(chan []byte)
	subAccountChan := make(chan []byte)
	//channel := make(chan []byte, 2)
	//var walletByte []byte
	//var subAccountByte []byte
	var walleturl = "https://api-v2-sandbox.chimoney.io/v0.2/multicurrency-wallets/create"
	//var wg sync.WaitGroup
	var subAccounturl = "https://api-v2-sandbox.chimoney.io/v0.2/sub-account/create"
	var walletConvertedBody map[string]interface{}
	var subAccountConvertedBody map[string]interface{}
	randomPhoneNumber := utils.RandomPhoneNumber()
	jsonResp, err := json.Marshal(types.NewChimoneySubAccount{
		Email:       payload.Email,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Name:        payload.FirstName + payload.LastName,
		PhoneNumber: randomPhoneNumber,
	})
	if err != nil {
		log.Printf("Error: %s", err)
		return nil, err
	}

	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	go func() {
		//Create a chimoney account
		request, _ := http.NewRequest("POST", subAccounturl, bytes.NewBuffer(jsonResp))
		request.Header.Add("content-type", "application/json")
		request.Header.Add("accept", "application/json")
		request.Header.Add("X-API-KEY", os.Getenv("CHIMONEY_API_KEY"))

		res, _ := client.Do(request)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		subAccountChan <- body
	}()

	go func() {
		//Create a chimoney account
		request, _ := http.NewRequest("POST", walleturl, bytes.NewBuffer(jsonResp))
		request.Header.Add("content-type", "application/json")
		request.Header.Add("accept", "application/json")
		request.Header.Add("X-API-KEY", os.Getenv("CHIMONEY_API_KEY"))

		res, _ := client.Do(request)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		walletChan <- body
	}()

	if err = json.Unmarshal(<-subAccountChan, &subAccountConvertedBody); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(<-walletChan, &walletConvertedBody); err != nil {
		return nil, err
	}

	defer close(subAccountChan)
	defer close(walletChan)

	//fmt.Println(string(<-walletChan))
	walletData := walletConvertedBody["data"].(map[string]interface{})
	subAccountData := subAccountConvertedBody["data"].(map[string]interface{})

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		log.Printf("Couldn't hash a password: %s", err)
		return nil, err
	}

	query := `INSERT INTO "user" (email, first_name, last_name, profile_picture, password, email_verified) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, email, first_name, last_name, profile_picture, email_verified`

	err = s.db.QueryRowContext(ctx, query, payload.Email, payload.FirstName, payload.LastName, payload.ProfilePicture, hashPassword, payload.EmailVerified).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.ProfilePicture,
		&user.EmailVerified,
	)

	query2 := `INSERT INTO "account_number" (subaccount_id, subaccount_number, user_id, wallet_id, wallet_number) VALUES ($1, $2, $3, $4, $5) RETURNING subaccount_number, wallet_number`
	err = s.db.QueryRowContext(ctx, query2, subAccountData["id"], subAccountData["parent"], user.ID, walletData["id"], walletData["parent"]).Scan(
		&user.SubAccountNumber,
		&user.WalletNumber,
	)
	//_, err = scanRowsToReturnUser(rows)
	return user, err
}

func scanRowsToReturnUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.EmailVerified,
		&user.ProfilePicture,
		&user.LastName,
		&user.FirstName,
		&user.ID,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GoogleAuthLoginAndRegister() error {
	return nil
}

//func (s *Store)
