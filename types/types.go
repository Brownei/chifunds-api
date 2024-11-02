package types

import (
	"context"
	"time"
)

type SSEChannel struct {
	Notifier chan string
	Clients  []chan string
}

type Balance struct {
	Amount int `json:"amount"`
}

type EncryptedDataPayload struct {
	Encrypted string `json:"encrypted"`
	AesKey    string `json:"aesKey"`
}

type User struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Password       string `json:"password"`
	EmailVerified  bool   `json:"email_verified"`
	AccountNumber  string `json:"account_number"`
	Balance        int32  `json:"balance"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=30"`
}

type RegisterUserPayload struct {
	Email          string `json:"email" validate:"required,email"`
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	ProfilePicture string `json:"profile_picture"`
	Password       string `json:"password" validate:"required,min=3,max=30"`
	EmailVerified  bool   `json:"email_verified"`
}

type BorrowMoneyDto struct {
	Explanation string `json:"explanation" validate:"required"`
	Amount      int32  `json:"amount" validate:"required"`
}

type TransferMoneyDto struct {
	Amount        int32  `json:"amount" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
}

type ReceivedTransactions struct {
	Amount          int32     `json:"amount"`
	SenderFirstName string    `json:"sender_first_name"`
	SenderLastName  string    `json:"sender_last_name"`
	SentAt          time.Time `json:"sent_at"`
}

type BorrowedTransactions struct {
	Amount int32     `json:"amount"`
	SentAt time.Time `json:"sent_at"`
}

type SentTransactions struct {
	Amount            int32     `json:"amount"`
	ReceiverFirstName string    `json:"receiver_first_name"`
	ReceiverLastName  string    `json:"receiver_last_name"`
	SentAt            time.Time `json:"sent_at"`
}

type DataPayloadFromTransferDto struct {
	Id             int8   `json:"id"`
	RemainingMoney int32  `json:"remaining_money"`
	Amount         int32  `json:"amount"`
	NameOfReceiver string `json:"name_of_receiver"`
}

type TransferJob struct {
	Id           int
	Query        string
	Args         []interface{}
	ExecuteQuery func(context.Context, string, []interface{}) error
}

type DataPayload struct {
	Data string `json:"data"`
}
