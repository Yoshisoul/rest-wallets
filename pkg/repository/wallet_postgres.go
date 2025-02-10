package repository

import (
	"fmt"
	"time"

	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WalletPostgres struct {
	db *sqlx.DB
}

func NewWalletPostgres(db *sqlx.DB) *WalletPostgres {
	return &WalletPostgres{db: db}
}

func (r *WalletPostgres) Create(userId int) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (wallet_id, user_id, amount, created_at, updated_at) values ($1, $2, $3, $4, $5) RETURNING wallet_id", walletTable)

	row := r.db.QueryRow(query, uuid.New(), userId, 0, time.Now(), time.Now())
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *WalletPostgres) GetAllFromUser(userId int) ([]models.Wallet, error) {
	var wallets []models.Wallet
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", walletTable)
	err := r.db.Select(&wallets, query, userId)

	return wallets, err
}

func (r *WalletPostgres) GetByIdFromUser(userId int, walletId uuid.UUID) (models.Wallet, error) {
	var wallet models.Wallet
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 AND wallet_id=$2", walletTable)
	err := r.db.Get(&wallet, query, userId, walletId)

	return wallet, err
}

func (r *WalletPostgres) GetById(walletId uuid.UUID) (models.Wallet, error) {
	var wallet models.Wallet
	query := fmt.Sprintf("SELECT * FROM %s WHERE wallet_id=$1", walletTable)
	err := r.db.Get(&wallet, query, walletId)

	return wallet, err
}

func (r *WalletPostgres) Delete(userId int, walletId uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	lockQuery := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 AND wallet_id=$2 FOR UPDATE", walletTable)
	_, err = tx.Exec(lockQuery, userId, walletId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1 AND wallet_id=$2", walletTable)
	_, err = tx.Exec(query, userId, walletId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
