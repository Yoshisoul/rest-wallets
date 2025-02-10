package service

import (
	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/Yoshisoul/rest-wallets/pkg/repository"
	"github.com/google/uuid"
)

type WalletService struct {
	repo repository.Wallet
}

func NewWalletService(repo repository.Wallet) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) Create(userId int) (uuid.UUID, error) {
	return s.repo.Create(userId)
}

func (s *WalletService) GetAllFromUser(userId int) ([]models.Wallet, error) {
	return s.repo.GetAllFromUser(userId)
}

func (s *WalletService) GetByIdFromUser(userId int, walletId uuid.UUID) (models.Wallet, error) {
	return s.repo.GetByIdFromUser(userId, walletId)
}

func (s *WalletService) Delete(userId int, walletId uuid.UUID) error {
	return s.repo.Delete(userId, walletId)
}
