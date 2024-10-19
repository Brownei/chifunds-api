package types

import "context"

type UserStore interface {
	GetUsersByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers() ([]User, error)
	CreateNewUser(ctx context.Context, payload RegisterUserPayload) (*User, error)
}

type AuthStore interface {
	GoogleAuthLoginAndRegister() error
}

type User struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	//Password       string `json:"password"`
	EmailVerified     bool     `json:"email_verified"`
	SubAccountNumber  string   `json:"subaccount_number"`
	SubAccountNumbers []string `json:"subaccount_numbers"`
	WalletNumber      string   `json:"wallet_number"`
	WalletNumbers     []string `json:"wallet_numbers"`

	//CreatedAt      time.Time `json:"created_at"`
}

type RegisterUserPayload struct {
	Email          string `json:"email" validate:"required,email"`
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	ProfilePicture string `json:"profile_picture"`
	Password       string `json:"password" validate:"required,min=3,max=30"`
	EmailVerified  bool   `json:"email_verified"`
}

type NewChimoneySubAccount struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}
