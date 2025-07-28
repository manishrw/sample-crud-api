package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/manishrw/sample-crud-api/go-echo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserServiceInterface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUser(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUsers() ([]*models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(id string, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	tests := []struct {
		name           string
		requestBody    models.CreateUserRequest
		mockReturn     *models.User
		mockError      error
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid user creation",
			requestBody: models.CreateUserRequest{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			mockReturn: &models.User{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "Missing required fields",
			requestBody: models.CreateUserRequest{
				Name:  "",
				Email: "john.doe@example.com",
			},
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Setup mock expectations only for successful cases
			if !tt.expectedError && tt.mockReturn != nil {
				mockService.On("CreateUser", &tt.requestBody).Return(tt.mockReturn, tt.mockError).Once()
			}

			// Test
			err := handler.CreateUser(c)

			// Assertions
			if tt.expectedError {
				assert.NoError(t, err) // Handler returns JSON, not error
				assert.Equal(t, tt.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetUsers(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	// Mock data
	mockUsers := []*models.User{
		{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		{
			Name:  "Jane Smith",
			Email: "jane.smith@example.com",
		},
	}

	// Setup mock expectations
	mockService.On("GetUsers").Return(mockUsers, nil).Once()

	// Setup request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err := handler.GetUsers(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockService.AssertExpectations(t)
}
