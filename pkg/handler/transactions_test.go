package handler

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/Yoshisoul/rest-wallets/pkg/service"
	mockService "github.com/Yoshisoul/rest-wallets/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_createTransaction(t *testing.T) {
	type mockBehavior func(s *mockService.MockTransaction, input models.TransactionInput)

	testTable := []struct {
		name                string
		inputBody           string
		mockExpInput        models.TransactionInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "Ok Deposit",
			inputBody: `{"walletId": "123e4567-e89b-12d3-a456-426614174000", "operationType":"DEPOSIT", "amount": 100}`,
			mockExpInput: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Deposit,
				Amount:        100,
			},
			mockBehavior: func(s *mockService.MockTransaction, input models.TransactionInput) {
				s.EXPECT().Create(input).Return(uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"uuid":"111e2222-e89b-12d3-a456-426614174000"}`,
		},
		{
			name:      "Ok Withdraw",
			inputBody: `{"walletId": "123e4567-e89b-12d3-a456-426614174000", "operationType":"WITHDRAW", "amount": 100}`,
			mockExpInput: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Withdraw,
				Amount:        100,
			},
			mockBehavior: func(s *mockService.MockTransaction, input models.TransactionInput) {
				s.EXPECT().Create(input).Return(uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"uuid":"111e2222-e89b-12d3-a456-426614174000"}`,
		},
		{
			name:                "Empty fields",
			inputBody:           `{"walletId": "123e4567-e89b-12d3-a456-426614174000", "operationType":"WITHDRAW"}`,
			mockBehavior:        func(s *mockService.MockTransaction, input models.TransactionInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:                "Incorrect fields",
			inputBody:           `{"walletId": "123e4567-e89b-12d3-a456-426614174000", "operationType":"INCORRECT", "amount": "100"}`,
			mockBehavior:        func(s *mockService.MockTransaction, input models.TransactionInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"walletId": "123e4567-e89b-12d3-a456-426614174000", "operationType":"DEPOSIT", "amount": 100}`,
			mockExpInput: models.TransactionInput{
				WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				OperationType: models.Deposit,
				Amount:        100,
			},
			mockBehavior: func(s *mockService.MockTransaction, input models.TransactionInput) {
				s.EXPECT().Create(input).Return(uuid.UUID{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transaction := mockService.NewMockTransaction(ctrl)
			testCase.mockBehavior(transaction, testCase.mockExpInput)

			services := &service.Service{Transaction: transaction}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/transactions", handler.createTransaction)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/transactions", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getAllTransactions(t *testing.T) {
	type mockBehavior func(s *mockService.MockTransaction)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mockService.MockTransaction) {
				s.EXPECT().GetAll().Return([]models.Transaction{
					{
						TransactionId: uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"),
						WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						OperationType: models.Deposit,
						Amount:        100,
						CreatedAt:     time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
					},
				}, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: `{"data":[{
			"transactionId":"111e2222-e89b-12d3-a456-426614174000",
			"walletId":"123e4567-e89b-12d3-a456-426614174000",
			"operationType":"DEPOSIT",
			"amount":100,
			"createdAt":"2025-02-10T00:00:00Z"}]}`,
		},
		{
			name: "Service Failure",
			mockBehavior: func(s *mockService.MockTransaction) {
				s.EXPECT().GetAll().Return([]models.Transaction{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name: "Empty",
			mockBehavior: func(s *mockService.MockTransaction) {
				s.EXPECT().GetAll().Return([]models.Transaction{}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[]}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transaction := mockService.NewMockTransaction(ctrl)
			testCase.mockBehavior(transaction)

			services := &service.Service{Transaction: transaction}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.GET("/transactions", handler.getAllTransactions)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/transactions", nil)

			// Perform Request
			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getTransactionById(t *testing.T) {
	type mockBehavior func(s *mockService.MockTransaction, id uuid.UUID)

	testTable := []struct {
		name                string
		inputId             string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:    "Ok",
			inputId: "111e2222-e89b-12d3-a456-426614174000",
			mockBehavior: func(s *mockService.MockTransaction, id uuid.UUID) {
				s.EXPECT().GetById(id).Return(models.Transaction{
					TransactionId: uuid.MustParse("111e2222-e89b-12d3-a456-426614174000"),
					WalletId:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					OperationType: models.Deposit,
					Amount:        100,
					CreatedAt:     time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: `{
			"transactionId":"111e2222-e89b-12d3-a456-426614174000",
			"walletId":"123e4567-e89b-12d3-a456-426614174000",
			"operationType":"DEPOSIT",
			"amount":100,
			"createdAt":"2025-02-10T00:00:00Z"}`,
		},
		{
			name:    "Service Failure",
			inputId: "111e2222-e89b-12d3-a456-426614174000",
			mockBehavior: func(s *mockService.MockTransaction, id uuid.UUID) {
				s.EXPECT().GetById(id).Return(models.Transaction{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:                "Invalid Transaction ID",
			inputId:             "invalid",
			mockBehavior:        func(s *mockService.MockTransaction, id uuid.UUID) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid id param"}`,
		},
		{
			name:    "Not found",
			inputId: "111e2222-e89b-12d3-a456-426614174000",
			mockBehavior: func(s *mockService.MockTransaction, id uuid.UUID) {
				s.EXPECT().GetById(id).Return(models.Transaction{}, sql.ErrNoRows)
			},
			expectedStatusCode:  404,
			expectedRequestBody: `{"message":"transaction not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			transaction := mockService.NewMockTransaction(c)
			if testCase.name != "Invalid Transaction ID" {
				testCase.mockBehavior(transaction, uuid.MustParse(testCase.inputId))
			}

			services := &service.Service{Transaction: transaction}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.GET("/transactions/:id", handler.getTransactionById)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/transactions/"+testCase.inputId, nil)

			// Perform Request
			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
