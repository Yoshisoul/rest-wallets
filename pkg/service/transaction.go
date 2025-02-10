package service

import (
	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/Yoshisoul/rest-wallets/pkg/repository"
	"github.com/google/uuid"
)

type TransactionService struct {
	repo       repository.Transaction
	walletRepo repository.Wallet
}

func NewTransactionService(repo repository.Transaction, walletRepo repository.Wallet) *TransactionService {
	return &TransactionService{repo: repo, walletRepo: walletRepo}
}

func (s *TransactionService) Create(transaction models.TransactionInput) (uuid.UUID, error) {
	_, err := s.walletRepo.GetById(transaction.WalletId)
	if err != nil {
		return uuid.Nil, err
	}

	return s.repo.Create(transaction)
}

func (s *TransactionService) GetAll() ([]models.Transaction, error) {
	return s.repo.GetAll()
}

func (s *TransactionService) GetById(transactionId uuid.UUID) (models.Transaction, error) {
	return s.repo.GetById(transactionId)
}
