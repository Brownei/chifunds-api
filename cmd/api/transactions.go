package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func (a *application) AllTransactionRoutes(r chi.Router) {
	r.Use(a.AuthMiddleware)
	r.Post("/transfer-money", a.TransferFunds)
	r.Post("/borrow-money", a.BorrowMoneyFromUs)
	r.Get("/received", a.GetReceivedTransactions)
	r.Get("/sent", a.GetSentTransactions)
}

func (a *application) BorrowMoneyFromUs(w http.ResponseWriter, r *http.Request) {
	var payload types.BorrowMoneyDto
	ctx := r.Context()
	email := ctx.Value("user").(string)

	decryptedData, err := utils.DecryptAndParseJson(r, RsaDecrypt)
	if err != nil {
		a.logger.Errorf("DECRYT ERROR: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := json.Unmarshal(decryptedData, &payload); err != nil {
		a.logger.Errorf("Unmarshall data error: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		fmt.Printf("Error: %s", errors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", errors))
		return
	}

	existingUser, _ := a.store.Users.GetUsersByEmail(ctx, email, false)
	if existingUser == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("There is no user like this"))
		return
	}

	if err := a.store.Transactions.BorrowMoney(ctx, payload.Amount, int8(existingUser.ID)); err != nil {
		a.logger.Info(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.EncryptAndWriteJson(w, http.StatusAccepted, []byte("Successfully"), RsaEncrypt)
	//utils.WriteJSON(w, http.StatusAccepted, "Successful")
}

func (a *application) TransferFunds(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUserEmail := ctx.Value("user").(string)

	decryptedData, err := utils.DecryptAndParseJson(r, RsaDecrypt)
	if err != nil {
		a.logger.Errorf("DECRYT ERROR: %v", err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	var payload types.TransferMoneyDto
	if err := json.Unmarshal(decryptedData, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	a.logger.Info(payload)
	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload: %v", errors))
		return
	}

	existingUser, _ := a.store.Users.GetUsersByEmail(ctx, currentUserEmail, false)
	if existingUser == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("There is no user like this"))
		return
	}

	if err := a.store.Transactions.TransferMoney(ctx, a.logger, *existingUser, payload.Amount, payload.AccountNumber); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.EncryptAndWriteJson(w, http.StatusAccepted, []byte("Successfully sent money"), RsaEncrypt)
	//utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Successfully sent money"))
}

func (a *application) GetReceivedTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := ctx.Value("email").(string)

	transactions, err := a.store.Transactions.GetReceivedTransactions(ctx, email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	byteTransactions, err := json.Marshal(transactions)

	utils.EncryptAndWriteJson(w, http.StatusOK, byteTransactions, RsaEncrypt)
}

func (a *application) GetSentTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := ctx.Value("email").(string)

	transactions, err := a.store.Transactions.GetSentTransactions(ctx, email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	byteTransactions, err := json.Marshal(transactions)

	utils.EncryptAndWriteJson(w, http.StatusOK, byteTransactions, RsaEncrypt)
}
