package repository

import (
	"fmt"
	"time"

	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) Create(transaction models.TransactionInput) (uuid.UUID, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (transaction_id, wallet_id, operation_type, amount, created_at) values ($1, $2, $3, $4, $5) RETURNING transaction_id", transactionTable)

	row := r.db.QueryRow(query, uuid.New(), transaction.WalletId, transaction.OperationType, transaction.Amount, time.Now())
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	lockQuery := fmt.Sprintf("SELECT amount FROM %s WHERE wallet_id = $1 FOR UPDATE", walletTable)
	_, err = tx.Exec(lockQuery, transaction.WalletId)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	var updateQuery string
	switch transaction.OperationType {
	case models.Deposit:
		updateQuery = fmt.Sprintf("UPDATE %s SET amount = amount + $1, updated_at = $2 WHERE wallet_id = $3", walletTable)
	case models.Withdraw:
		updateQuery = fmt.Sprintf("UPDATE %s SET amount = amount - $1, updated_at = $2 WHERE wallet_id = $3", walletTable)
	}

	_, err = tx.Exec(updateQuery, transaction.Amount, time.Now(), transaction.WalletId)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *TransactionPostgres) GetAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := fmt.Sprintf("SELECT * FROM %s", transactionTable)
	err := r.db.Select(&transactions, query)

	return transactions, err
}

func (r *TransactionPostgres) GetById(transactionId uuid.UUID) (models.Transaction, error) {
	var transaction models.Transaction
	query := fmt.Sprintf("SELECT * FROM %s WHERE transaction_id=$1", transactionTable)
	err := r.db.Get(&transaction, query, transactionId)

	return transaction, err
}
