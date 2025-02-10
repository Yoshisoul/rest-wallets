package service

import (
	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/Yoshisoul/rest-wallets/pkg/repository"
	"github.com/google/uuid"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user models.SignUpInput) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Wallet interface {
	Create(userId int) (uuid.UUID, error)
	GetAllFromUser(userId int) ([]models.Wallet, error)
	GetByIdFromUser(userId int, walletId uuid.UUID) (models.Wallet, error)
	Delete(userId int, walletId uuid.UUID) error
}

type Transaction interface {
	Create(transaction models.TransactionInput) (uuid.UUID, error)
	GetAll() ([]models.Transaction, error)
	GetById(transactionId uuid.UUID) (models.Transaction, error)
}

type Service struct {
	Authorization
	Wallet
	Transaction
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Wallet:        NewWalletService(repos.Wallet),
		Transaction:   NewTransactionService(repos.Transaction, repos.Wallet),
	}
}
