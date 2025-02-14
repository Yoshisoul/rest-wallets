package handler

import (
	"database/sql"
	"errors"
	"net/http/httptest"
	"strconv"
	"strings"
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

func setUserIdMiddleware(userId int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if userId != -1 { // to simulate a missing user id
			c.Set(userCtx, userId)
		}
		c.Next()
	}
}

func TestHandler_createWallet(t *testing.T) {
	type mockBehavior func(s *mockService.MockWallet, id int)

	testTable := []struct {
		name                string
		inputUserId         int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			inputUserId: 1,
			mockBehavior: func(s *mockService.MockWallet, id int) {
				s.EXPECT().Create(id).Return(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"uuid":"123e4567-e89b-12d3-a456-426614174000"}`,
		},
		{
			name:        "Service Failure",
			inputUserId: 1,
			mockBehavior: func(s *mockService.MockWallet, id int) {
				s.EXPECT().Create(id).Return(uuid.UUID{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:                "UserID not found",
			inputUserId:         -1,
			mockBehavior:        func(s *mockService.MockWallet, id int) {},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			wallet := mockService.NewMockWallet(c)
			testCase.mockBehavior(wallet, testCase.inputUserId)

			services := &service.Service{Wallet: wallet}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.Use(setUserIdMiddleware(testCase.inputUserId))
			r.POST("/wallets", handler.createWallet)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/wallets", strings.NewReader(strconv.Itoa(testCase.inputUserId)))
			req.Header.Set("Authorization", "Bearer token")

			// Perform Request
			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getAllWalletsFromUser(t *testing.T) {
	type mockBehavior func(s *mockService.MockWallet, id int)

	testTable := []struct {
		name                string
		inputUserId         int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			inputUserId: 1,
			mockBehavior: func(s *mockService.MockWallet, id int) {
				s.EXPECT().GetAllFromUser(id).Return([]models.Wallet{
					{
						WalletId:  uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
						UserId:    id,
						Amount:    100,
						CreatedAt: time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
					},
				}, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: `{"data":[{
			"walletId":"123e4567-e89b-12d3-a456-426614174000",
			"userId":1,
			"amount":100,
			"createdAt":"2025-02-10T00:00:00Z",
			"updatedAt":"2025-02-10T00:00:00Z"}]}`,
		},
		{
			name:        "Service Failure",
			inputUserId: 1,
			mockBehavior: func(s *mockService.MockWallet, id int) {
				s.EXPECT().GetAllFromUser(id).Return(nil, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:                "UserID not found",
			inputUserId:         -1,
			mockBehavior:        func(s *mockService.MockWallet, id int) {},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:        "Empty",
			inputUserId: 1,
			mockBehavior: func(s *mockService.MockWallet, id int) {
				s.EXPECT().GetAllFromUser(id).Return([]models.Wallet{}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[]}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			wallet := mockService.NewMockWallet(c)
			testCase.mockBehavior(wallet, testCase.inputUserId)

			services := &service.Service{Wallet: wallet}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.Use(setUserIdMiddleware(testCase.inputUserId))
			r.GET("/wallets", handler.getAllWalletsFromUser)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/wallets", nil)
			req.Header.Set("Authorization", "Bearer token")

			// Perform Request
			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getWalletById(t *testing.T) {
	type mockBehavior func(s *mockService.MockWallet, userId int, walletId uuid.UUID)

	testTable := []struct {
		name                string
		inputUserId         int
		inputWalletId       string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:          "OK",
			inputUserId:   1,
			inputWalletId: "123e4567-e89b-12d3-a456-426614174000",
			mockBehavior: func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {
				s.EXPECT().GetByIdFromUser(userId, walletId).Return(models.Wallet{
					WalletId:  walletId,
					UserId:    userId,
					Amount:    100,
					CreatedAt: time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: `{"walletId":"123e4567-e89b-12d3-a456-426614174000",
			"userId":1,
			"amount":100,
			"createdAt":"2025-02-10T00:00:00Z",
			"updatedAt":"2025-02-10T00:00:00Z"}`,
		},
		{
			name:          "Service Failure",
			inputUserId:   1,
			inputWalletId: "123e4567-e89b-12d3-a456-426614174000",
			mockBehavior: func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {
				s.EXPECT().GetByIdFromUser(userId, walletId).Return(models.Wallet{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:                "UserID not found",
			inputUserId:         -1,
			inputWalletId:       "123e4567-e89b-12d3-a456-426614174000",
			mockBehavior:        func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:                "Invalid Wallet ID",
			inputUserId:         1,
			inputWalletId:       "invalid",
			mockBehavior:        func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid id param"}`,
		},
		{
			name:          "Wallet not found",
			inputUserId:   1,
			inputWalletId: "123e4567-e89b-12d3-a456-426614174123",
			mockBehavior: func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {
				s.EXPECT().GetByIdFromUser(userId, walletId).Return(models.Wallet{}, sql.ErrNoRows)
			},
			expectedStatusCode:  404,
			expectedRequestBody: `{"message":"wallet not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			wallet := mockService.NewMockWallet(c)
			if testCase.name != "Invalid Wallet ID" {
				testCase.mockBehavior(wallet, testCase.inputUserId, uuid.MustParse(testCase.inputWalletId))
			}

			services := &service.Service{Wallet: wallet}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.Use(setUserIdMiddleware(testCase.inputUserId))
			r.GET("/wallets/:id", handler.getWalletById)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/wallets/"+testCase.inputWalletId, nil)
			req.Header.Set("Authorization", "Bearer token")

			// Perform Request
			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_deleteWallet(t *testing.T) {
	type mockBehavior func(s *mockService.MockWallet, userId int, walletId uuid.UUID)

	testTable := []struct {
		name                string
		inputUserId         int
		inputWalletId       string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:          "OK",
			inputUserId:   1,
			inputWalletId: "123e4567-e89b-12d3-a456-426614174000",
			mockBehavior: func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {
				s.EXPECT().Delete(userId, walletId).Return(nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"status":"ok"}`,
		},
		{
			name:          "Service Failure",
			inputUserId:   1,
			inputWalletId: "123e4567-e89b-12d3-a456-426614174000",
			mockBehavior: func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {
				s.EXPECT().Delete(userId, walletId).Return(errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
		{
			name:                "UserID not found",
			inputUserId:         -1,
			inputWalletId:       "123e4567-e89b-12d3-a456-426614174000",
			mockBehavior:        func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"user id not found"}`,
		},
		{
			name:                "Invalid Wallet ID",
			inputUserId:         1,
			inputWalletId:       "invalid",
			mockBehavior:        func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid id param"}`,
		},
		{
			name:          "Wallet not found",
			inputUserId:   1,
			inputWalletId: "123e4567-e89b-12d3-a456-426614174123",
			mockBehavior: func(s *mockService.MockWallet, userId int, walletId uuid.UUID) {
				s.EXPECT().Delete(userId, walletId).Return(errors.New("wallet not found"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"wallet not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			wallet := mockService.NewMockWallet(c)
			if testCase.name != "Invalid Wallet ID" {
				testCase.mockBehavior(wallet, testCase.inputUserId, uuid.MustParse(testCase.inputWalletId))
			}

			services := &service.Service{Wallet: wallet}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.Use(setUserIdMiddleware(testCase.inputUserId))
			r.DELETE("/wallets/:id", handler.deleteWallet)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/wallets/"+testCase.inputWalletId, nil)
			req.Header.Set("Authorization", "Bearer token")

			// Perform Request
			r.ServeHTTP(w, req)

			// Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
