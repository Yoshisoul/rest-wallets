package repository

import (
	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.SignUpInput) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Wallet interface {
	Create(userId int) (uuid.UUID, error)
	GetAllFromUser(userId int) ([]models.Wallet, error)
	GetByIdFromUser(userId int, walletId uuid.UUID) (models.Wallet, error)
	GetById(walletId uuid.UUID) (models.Wallet, error)
	Delete(userId int, walletId uuid.UUID) error
}

type Transaction interface {
	Create(transaction models.TransactionInput) (uuid.UUID, error)
	GetAll() ([]models.Transaction, error)
	GetById(transactionId uuid.UUID) (models.Transaction, error)
}

type Repository struct {
	Authorization
	Wallet
	Transaction
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Wallet:        NewWalletPostgres(db),
		Transaction:   NewTransactionPostgres(db),
	}
}
