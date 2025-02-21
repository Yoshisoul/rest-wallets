package models

import (
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type Wallet struct {
	WalletId  uuid.UUID `json:"walletId" db:"wallet_id"`
	UserId    int       `json:"userId" db:"user_id"`
	Amount    int64     `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Transaction struct {
	TransactionId uuid.UUID     `json:"transactionId" db:"transaction_id"`
	WalletId      uuid.UUID     `json:"walletId" db:"wallet_id" binding:"required"`
	OperationType OperationType `json:"operationType" db:"operation_type" binding:"required"`
	Amount        int64         `json:"amount" db:"amount" binding:"required"`
	CreatedAt     time.Time     `json:"createdAt" db:"created_at"`
}

type TransactionInput struct {
	WalletId      uuid.UUID     `json:"walletId" db:"wallet_id" binding:"required"`
	OperationType OperationType `json:"operationType" db:"operation_type" binding:"required"`
	Amount        int64         `json:"amount" db:"amount" binding:"required"`
}
