package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/Yoshisoul/rest-wallets/pkg/service"
	mockService "github.com/Yoshisoul/rest-wallets/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_signup(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, user models.SignUpInput)

	testTable := []struct {
		name                string
		inputBody           string
		mockExpInput        models.SignUpInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "test", "username":"user", "password": "pass"}`,
			mockExpInput: models.SignUpInput{
				Name:     "test",
				Username: "user",
				Password: "pass",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user models.SignUpInput) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty fields",
			inputBody:           `{"username":"user", "password": "pass"}`,
			mockBehavior:        func(s *mockService.MockAuthorization, user models.SignUpInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"name": "test", "username":"user", "password": "pass"}`,
			mockExpInput: models.SignUpInput{
				Name:     "test",
				Username: "user",
				Password: "pass",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user models.SignUpInput) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.mockExpInput)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, user models.SignInInput)

	testTable := []struct {
		testName            string
		inputBody           string
		mockExpInput        models.SignInInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			testName:  "OK",
			inputBody: `{"username":"user", "password": "pass"}`,
			mockExpInput: models.SignInInput{
				Username: "user",
				Password: "pass",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user models.SignInInput) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("token", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"token":"token"}`,
		},
		{
			testName:            "Empty fields",
			inputBody:           `{"password": "pass"}`,
			mockBehavior:        func(s *mockService.MockAuthorization, user models.SignInInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			testName:  "Service Failure",
			inputBody: `{"username":"user", "password": "pass"}`,
			mockExpInput: models.SignInInput{
				Username: "user",
				Password: "pass",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user models.SignInInput) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("", errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.testName, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.mockExpInput)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Test Server
			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-in",
				bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
