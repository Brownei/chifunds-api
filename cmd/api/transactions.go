package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	r.Get("/borrowed", a.GetBorrowedTransactions)
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

	if err := a.store.Transactions.BorrowMoney(ctx, a.logger, payload.Amount, int8(existingUser.ID)); err != nil {
		a.logger.Info(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	balance, err := a.store.Users.GetBalance(ctx, existingUser.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	broadcastBalanceUpdate(strconv.Itoa(balance.Amount))

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

	balance, err := a.store.Users.GetBalance(ctx, existingUser.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	broadcastBalanceUpdate(strconv.Itoa(balance.Amount))

	utils.EncryptAndWriteJson(w, http.StatusAccepted, []byte("Successfully sent money"), RsaEncrypt)
	//utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("Successfully sent money"))
}

func (a *application) GetReceivedTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := ctx.Value("user").(string)

	transactions, err := a.store.Transactions.GetReceivedTransactions(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("No transactions received"))
			return
		}
		a.logger.Info(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, transactions)
}

func (a *application) GetBorrowedTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	transactions, err := a.store.Transactions.GetBorrowedTransactions(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("You do not want to borow money yet"))
			return
		}
		a.logger.Info(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	a.logger.Infof("Transactions Borrowed: %v", transactions)

	utils.WriteJSON(w, http.StatusOK, transactions)
}

func (a *application) GetSentTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	email := ctx.Value("user").(string)

	transactions, err := a.store.Transactions.GetSentTransactions(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("No transactions sent"))
			return
		}
		a.logger.Info(err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, transactions)
}
