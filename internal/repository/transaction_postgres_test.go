package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/Yoshisoul/rest-wallets/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTransaction_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	r := NewTransactionPostgres(db)

	type mockBehavior func(input models.TransactionInput)

	testTable := []struct {
		name         string
		input        models.TransactionInput
		expectedId   uuid.UUID
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "Ok Deposit",
			input: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Deposit,
				Amount:        100,
			},
			expectedId: uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"),
			mockBehavior: func(input models.TransactionInput) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs(sqlmock.AnyArg(), input.WalletId, input.OperationType, input.Amount, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}).
						AddRow("111e2222-e89b-12d3-a456-426614174000"))
				mock.ExpectExec("SELECT amount FROM wallets WHERE wallet_id = \\$1 FOR UPDATE").
					WithArgs(input.WalletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE wallets SET amount = amount \\+ \\$1, updated_at = \\$2 WHERE wallet_id = \\$3").
					WithArgs(input.Amount, sqlmock.AnyArg(), input.WalletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Ok Withdraw",
			input: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Withdraw,
				Amount:        100,
			},
			expectedId: uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"),
			mockBehavior: func(input models.TransactionInput) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs(sqlmock.AnyArg(), input.WalletId, input.OperationType, input.Amount, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}).
						AddRow("111e2222-e89b-12d3-a456-426614174000"))
				mock.ExpectExec("SELECT amount FROM wallets WHERE wallet_id = \\$1 FOR UPDATE").
					WithArgs(input.WalletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE wallets SET amount = amount \\- \\$1, updated_at = \\$2 WHERE wallet_id = \\$3").
					WithArgs(input.Amount, sqlmock.AnyArg(), input.WalletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Insert Error, rollback",
			input: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Deposit,
				Amount:        100,
			},
			mockBehavior: func(input models.TransactionInput) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs(sqlmock.AnyArg(), input.WalletId, input.OperationType, input.Amount, sqlmock.AnyArg()).
					WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Select Error, rollback",
			input: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Deposit,
				Amount:        100,
			},
			mockBehavior: func(input models.TransactionInput) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs(sqlmock.AnyArg(), input.WalletId, input.OperationType, input.Amount, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}).
						AddRow("111e2222-e89b-12d3-a456-426614174000"))
				mock.ExpectExec("SELECT amount FROM wallets WHERE wallet_id = \\$1 FOR UPDATE").
					WithArgs(input.WalletId).
					WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Update Error, rollback",
			input: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Deposit,
				Amount:        100,
			},
			mockBehavior: func(input models.TransactionInput) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs(sqlmock.AnyArg(), input.WalletId, input.OperationType, input.Amount, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}).
						AddRow("111e2222-e89b-12d3-a456-426614174000"))
				mock.ExpectExec("SELECT amount FROM wallets WHERE wallet_id = \\$1 FOR UPDATE").
					WithArgs(input.WalletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE wallets SET amount = amount \\+ \\$1, updated_at = \\$2 WHERE wallet_id = \\$3").
					WithArgs(input.Amount, sqlmock.AnyArg(), input.WalletId).
					WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Commit Error",
			input: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Deposit,
				Amount:        100,
			},
			mockBehavior: func(input models.TransactionInput) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO transactions").
					WithArgs(sqlmock.AnyArg(), input.WalletId, input.OperationType, input.Amount, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}).
						AddRow("111e2222-e89b-12d3-a456-426614174000"))
				mock.ExpectExec("SELECT amount FROM wallets WHERE wallet_id = \\$1 FOR UPDATE").
					WithArgs(input.WalletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE wallets SET amount = amount \\+ \\$1, updated_at = \\$2 WHERE wallet_id = \\$3").
					WithArgs(input.Amount, sqlmock.AnyArg(), input.WalletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "Begin Error, rollback",
			input: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Deposit,
				Amount:        100,
			},
			mockBehavior: func(input models.TransactionInput) {
				mock.ExpectBegin().WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.mockBehavior(testcase.input)
			got, err := r.Create(testcase.input)
			if testcase.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testcase.expectedId, got)
			}
		})
	}
}

func TestTransaction_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	r := NewTransactionPostgres(db)
	type mockBehavior func()

	testTable := []struct {
		name         string
		expected     []models.Transaction
		mockBehavior mockBehavior
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "Ok",
			expected: []models.Transaction{
				{
					TransactionId: uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"),
					WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					OperationType: models.Deposit,
					Amount:        100,
					CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			mockBehavior: func() {
				mock.ExpectQuery("SELECT \\* FROM transactions").
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "wallet_id", "operation_type", "amount", "created_at"}).
						AddRow("111e2222-e89b-12d3-a456-426614174000", "123e4567-e89b-12d3-a456-426614174000", models.Deposit, 100, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
			},
		},
		{
			name:     "Ok, empty",
			expected: []models.Transaction{},
			mockBehavior: func() {
				mock.ExpectQuery("SELECT \\* FROM transactions").
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "wallet_id", "operation_type", "amount", "created_at"}))
			},
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.mockBehavior()

			got, err := r.GetAll()
			if testcase.wantErr {
				assert.Error(t, err)
				assert.Equal(t, testcase.expectedErr, err)
				return
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, testcase.expected, got)
			}
		})
	}
}

func TestTransaction_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	r := NewTransactionPostgres(db)
	type mockBehavior func(id uuid.UUID)

	testTable := []struct {
		name         string
		inputId      uuid.UUID
		expected     models.Transaction
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name:    "Ok",
			inputId: uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"),
			expected: models.Transaction{
				TransactionId: uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"),
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Deposit,
				Amount:        100,
				CreatedAt:     time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			mockBehavior: func(id uuid.UUID) {
				mock.ExpectQuery("SELECT \\* FROM transactions WHERE transaction_id = \\$1").
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "wallet_id", "operation_type", "amount", "created_at"}).
						AddRow("111e2222-e89b-12d3-a456-426614174000", "123e4567-e89b-12d3-a456-426614174000", models.Deposit, 100, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
			},
		},
		{
			name:     "Not Found",
			inputId:  uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"),
			expected: models.Transaction{},
			mockBehavior: func(id uuid.UUID) {
				mock.ExpectQuery("SELECT \\* FROM transactions WHERE transaction_id = \\$1").
					WithArgs(id).
					WillReturnError(errors.New("sql: no rows in result set"))
			},
			wantErr: true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.mockBehavior(testcase.inputId)

			got, err := r.GetById(testcase.inputId)
			if testcase.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testcase.expected, got)
			}
		})
	}
}
