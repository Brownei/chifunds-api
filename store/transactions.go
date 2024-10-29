package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/types/worker"
	"go.uber.org/zap"
)

type TransactionStore struct {
	store *sql.DB
}

func (s *TransactionStore) BorrowMoney(ctx context.Context, lendedMoney int32, userId int8) error {
	var accountId int
	if err := s.store.QueryRowContext(
		ctx,
		`UPDATE "account" SET money = money + $1 WHERE user_id = $2 RETURNING id`,
		lendedMoney,
		userId,
	).Scan(
		&accountId,
	); err != nil {
		return err
	}
	return nil
}

func (s *TransactionStore) TransferMoney(ctx context.Context, logger *zap.SugaredLogger, existingUser types.User, amount int32, accountNumber string) error {
	transactionWorker := worker.NewWorker(ctx, logger)
	receivingUser, err := s.GetAccountFromAccountNumber(ctx, accountNumber)
	if err != nil {
		return err
	}

	if receivingUser == nil {
		return fmt.Errorf("No user with this account number")
	}

	if existingUser.Balance < amount {
		return fmt.Errorf("Insufficient funds")
	}

	jobs := []types.TransferJob{
		{
			Id:           1,
			Query:        `UPDATE "account" SET money = money - $1 WHERE user_id = $2 RETURNING id`,
			ExecuteQuery: s.ExecuteQuery,
			Args:         []interface{}{amount, existingUser.ID},
		},
		{
			Id:           2,
			Query:        `UPDATE "account" SET money = money + $1 WHERE user_id = $2 RETURNING id`,
			ExecuteQuery: s.ExecuteQuery,
			Args:         []interface{}{amount, receivingUser.ID},
		},
		{
			Id:           3,
			Query:        `INSERT INTO "transactions" (receiver_id, sender_id, amount_sent) VALUES ($1, $2, $3) RETURNING id`,
			ExecuteQuery: s.ExecuteQuery,
			Args:         []interface{}{receivingUser.ID, existingUser.ID, amount},
		},
	}

	transactionWorker.RunQueriesWithWorkerPool(jobs, len(jobs))
	return nil
}

func (s *TransactionStore) ExecuteQuery(ctx context.Context, query string, args []interface{}) error {
	var datapayload types.DataPayloadFromTransferDto
	if err := s.store.QueryRowContext(ctx, query, args...).Scan(&datapayload.Id); err != nil {
		return err
	}

	return nil
}

func (s *TransactionStore) GetAccountFromAccountNumber(ctx context.Context, accountNumber string) (*types.User, error) {
	user := &types.User{}
	query := `SELECT u.id, u.first_name, u.last_name, u.profile_picture, a.account_number FROM "user" AS u JOIN "account" AS a ON u.id = a.user_id WHERE a.account_number = $1`

	if err := s.store.QueryRowContext(ctx, query, accountNumber).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.ProfilePicture,
		&user.AccountNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No user with this account!")
		}
		return nil, err
	}

	return user, nil
}

func (s *TransactionStore) GetReceivedTransactions(ctx context.Context, email string) (*types.ReceivedTransactions, error) {
	transactions := &types.ReceivedTransactions{}
	query := `SELECT t.amount_sent, t.sent_at, r.first_name, r.last_name FROM "user" AS r JOIN "transactions" AS t ON t.receiver_id = r.id WHERE email = $1`

	if err := s.store.QueryRowContext(ctx, query, email).Scan(
		&transactions.Amount,
		&transactions.SentAt,
		&transactions.ReceiverFirstName,
		&transactions.ReceiverLastName,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No transactions received!")
		}
	}

	return transactions, nil
}

func (s *TransactionStore) GetSentTransactions(ctx context.Context, email string) (*types.SentTransactions, error) {
	transactions := &types.SentTransactions{}
	query := `SELECT t.amount_sent, t.sent_at, r.first_name, r.last_name FROM "user" AS r JOIN "transactions" AS t ON t.sender_id = r.id WHERE email = $1`

	if err := s.store.QueryRowContext(ctx, query, email).Scan(
		&transactions.Amount,
		&transactions.SentAt,
		&transactions.ReceiverFirstName,
		&transactions.ReceiverLastName,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No transactions received!")
		}
	}

	return transactions, nil
}
