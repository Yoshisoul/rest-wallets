package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestWallet_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewWalletPostgres(db)

	type mockBehavior func(userId int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		userId       int
		wantErr      bool
	}{
		{
			name:   "OK",
			userId: 1,
			mockBehavior: func(userId int) {
				rows := sqlmock.NewRows([]string{"wallet_id"}).AddRow(uuid.New())
				mock.ExpectQuery("INSERT INTO wallets").
					WithArgs(sqlmock.AnyArg(), userId, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
		},
		{
			name:   "Some error",
			userId: 1,
			mockBehavior: func(userId int) {
				mock.ExpectQuery("INSERT INTO wallets").
					WithArgs(sqlmock.AnyArg(), userId, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "Row error",
			mockBehavior: func(userId int) {
				rows := sqlmock.NewRows([]string{"wallet_id"}).
					AddRow(uuid.New()).
					RowError(0, errors.New("some row error"))
				mock.ExpectQuery("INSERT INTO wallets").
					WithArgs(sqlmock.AnyArg(), userId, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userId)

			got, err := r.Create(testCase.userId)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, got)
				assert.IsType(t, uuid.UUID{}, got)
			}
		})
	}
}

func TestWallet_GetAllFromUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewWalletPostgres(db)

	type mockBehavior func(userId int, expectedOut []models.Wallet)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		inputUserId  int
		expectedOut  []models.Wallet
		wantErr      bool
	}{
		{
			name:        "OK",
			inputUserId: 1,
			expectedOut: []models.Wallet{
				{WalletId: uuid.New(), UserId: 1, Amount: 50, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				{WalletId: uuid.New(), UserId: 1, Amount: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			},
			mockBehavior: func(userId int, expectedOut []models.Wallet) {
				rows := sqlmock.NewRows([]string{"wallet_id", "user_id", "amount", "created_at", "updated_at"}).
					AddRow(expectedOut[0].WalletId, expectedOut[0].UserId, expectedOut[0].Amount, expectedOut[0].CreatedAt, expectedOut[0].UpdatedAt).
					AddRow(expectedOut[1].WalletId, expectedOut[1].UserId, expectedOut[1].Amount, expectedOut[1].CreatedAt, expectedOut[1].UpdatedAt)
				mock.ExpectQuery(`SELECT \* FROM wallets WHERE user_id=\$1`).
					WithArgs(userId).
					WillReturnRows(rows)
			},
		},
		{
			name:        "Ok, empty result",
			inputUserId: 1,
			mockBehavior: func(userId int, expectedOut []models.Wallet) {
				rows := sqlmock.NewRows([]string{"wallet_id", "user_id", "amount", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT \* FROM wallets WHERE user_id=\$1`).
					WithArgs(userId).
					WillReturnRows(rows)
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.inputUserId, testCase.expectedOut)

			got, err := r.GetAllFromUser(testCase.inputUserId)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, testCase.expectedOut, got)
			}
		})
	}
}

func TestWallet_GetByIdFromUser(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewWalletPostgres(db)

	type mockBehavior func(userId int, walletId uuid.UUID, expectedOut models.Wallet)

	testTable := []struct {
		name          string
		mockBehavior  mockBehavior
		inputUserId   int
		inputWalletId uuid.UUID
		expectedOut   models.Wallet
		wantErr       bool
	}{
		{
			name:          "OK",
			inputUserId:   1,
			inputWalletId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expectedOut: models.Wallet{
				WalletId:  uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				UserId:    1,
				Amount:    50,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockBehavior: func(userId int, walletId uuid.UUID, expectedOut models.Wallet) {
				rows := sqlmock.NewRows([]string{"wallet_id", "user_id", "amount", "created_at", "updated_at"}).
					AddRow(expectedOut.WalletId, expectedOut.UserId, expectedOut.Amount, expectedOut.CreatedAt, expectedOut.UpdatedAt)
				mock.ExpectQuery(`SELECT \* FROM wallets WHERE user_id=\$1 AND wallet_id=\$2`).
					WithArgs(userId, walletId).
					WillReturnRows(rows)
			},
		},
		{
			name:          "Not found",
			inputUserId:   1,
			inputWalletId: uuid.New(),
			mockBehavior: func(userId int, walletId uuid.UUID, expectedOut models.Wallet) {
				rows := sqlmock.NewRows([]string{"wallet_id", "user_id", "amount", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT \* FROM wallets WHERE user_id=\$1 AND wallet_id=\$2`).
					WithArgs(userId, walletId).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.inputUserId, testCase.inputWalletId, testCase.expectedOut)

			got, err := r.GetByIdFromUser(testCase.inputUserId, testCase.inputWalletId)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedOut, got)
			}
		})
	}
}

func TestWallet_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewWalletPostgres(db)

	type mockBehavior func(walletId uuid.UUID, expectedOut models.Wallet)

	testTable := []struct {
		name          string
		mockBehavior  mockBehavior
		inputWalletId uuid.UUID
		expectedOut   models.Wallet
		wantErr       bool
	}{
		{
			name:          "OK",
			inputWalletId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			expectedOut: models.Wallet{
				WalletId:  uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				UserId:    1,
				Amount:    50,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockBehavior: func(walletId uuid.UUID, expectedOut models.Wallet) {
				rows := sqlmock.NewRows([]string{"wallet_id", "user_id", "amount", "created_at", "updated_at"}).
					AddRow(expectedOut.WalletId, expectedOut.UserId, expectedOut.Amount, expectedOut.CreatedAt, expectedOut.UpdatedAt)
				mock.ExpectQuery(`SELECT \* FROM wallets WHERE wallet_id=\$1`).
					WithArgs(walletId).
					WillReturnRows(rows)
			},
		},
		{
			name:          "Not found",
			inputWalletId: uuid.New(),
			mockBehavior: func(walletId uuid.UUID, expectedOut models.Wallet) {
				rows := sqlmock.NewRows([]string{"wallet_id", "user_id", "amount", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT \* FROM wallets WHERE wallet_id=\$1`).
					WithArgs(walletId).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.inputWalletId, testCase.expectedOut)

			got, err := r.GetById(testCase.inputWalletId)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedOut, got)
			}
		})
	}
}

func TestWallet_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewWalletPostgres(db)

	type mockBehavior func(userId int, walletId uuid.UUID)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		userId       int
		walletId     uuid.UUID
		wantErr      bool
	}{
		{
			name:     "OK",
			userId:   1,
			walletId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			mockBehavior: func(userId int, walletId uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec(`SELECT \* FROM wallets WHERE user_id=\$1 AND wallet_id=\$2 FOR UPDATE`).
					WithArgs(userId, walletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(`DELETE FROM wallets WHERE user_id=\$1 AND wallet_id=\$2`).
					WithArgs(userId, walletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:     "Begin error",
			userId:   1,
			walletId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			mockBehavior: func(userId int, walletId uuid.UUID) {
				mock.ExpectBegin().WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name:     "Lock error, rollback",
			userId:   1,
			walletId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			mockBehavior: func(userId int, walletId uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec(`SELECT \* FROM wallets WHERE user_id=\$1 AND wallet_id=\$2 FOR UPDATE`).
					WithArgs(userId, walletId).
					WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:     "Delete error, rollback",
			userId:   1,
			walletId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			mockBehavior: func(userId int, walletId uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec(`SELECT \* FROM wallets WHERE user_id=\$1 AND wallet_id=\$2 FOR UPDATE`).
					WithArgs(userId, walletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(`DELETE FROM wallets WHERE user_id=\$1 AND wallet_id=\$2`).
					WithArgs(userId, walletId).
					WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name:     "Commit error",
			userId:   1,
			walletId: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			mockBehavior: func(userId int, walletId uuid.UUID) {
				mock.ExpectBegin()
				mock.ExpectExec(`SELECT \* FROM wallets WHERE user_id=\$1 AND wallet_id=\$2 FOR UPDATE`).
					WithArgs(userId, walletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(`DELETE FROM wallets WHERE user_id=\$1 AND wallet_id=\$2`).
					WithArgs(userId, walletId).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.userId, testCase.walletId)

			err := r.Delete(testCase.userId, testCase.walletId)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
