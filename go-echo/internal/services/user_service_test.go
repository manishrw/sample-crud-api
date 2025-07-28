package services

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/manishrw/sample-crud-api/go-echo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]*models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		request        *models.CreateUserRequest
		mockSetup      func(*MockUserRepository)
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "Successfully create user",
			request: &models.CreateUserRequest{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByEmail", "john.doe@example.com").Return(nil, fmt.Errorf("user not found"))
				mockRepo.On("Create", mock.MatchedBy(func(user *models.User) bool {
					return user.Name == "John Doe" && user.Email == "john.doe@example.com"
				})).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "User with email already exists",
			request: &models.CreateUserRequest{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				existingUser := &models.User{
					ID:    uuid.New(),
					Name:  "Existing User",
					Email: "john.doe@example.com",
				}
				mockRepo.On("GetByEmail", "john.doe@example.com").Return(existingUser, nil)
			},
			expectedError:  true,
			expectedErrMsg: "user with email john.doe@example.com already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUserService(mockRepo)

			tt.mockSetup(mockRepo)

			user, err := service.CreateUser(tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.request.Name, user.Name)
				assert.Equal(t, tt.request.Email, user.Email)
				// Note: UUID is set by repository, not service
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserRepository)
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:   "Successfully get user",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(mockRepo *MockUserRepository) {
				userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				expectedUser := &models.User{
					ID:    userID,
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}
				mockRepo.On("GetByID", userID).Return(expectedUser, nil)
			},
			expectedError: false,
		},
		{
			name:   "Invalid UUID format",
			userID: "invalid-uuid",
			mockSetup: func(mockRepo *MockUserRepository) {
				// No mock setup needed for invalid UUID
			},
			expectedError:  true,
			expectedErrMsg: "invalid user ID format",
		},
		{
			name:   "User not found",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(mockRepo *MockUserRepository) {
				userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				mockRepo.On("GetByID", userID).Return(nil, fmt.Errorf("user not found"))
			},
			expectedError:  true,
			expectedErrMsg: "failed to get user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUserService(mockRepo)

			tt.mockSetup(mockRepo)

			user, err := service.GetUser(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUsers := []*models.User{
		{
			ID:    uuid.New(),
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		{
			ID:    uuid.New(),
			Name:  "Jane Smith",
			Email: "jane.smith@example.com",
		},
	}

	mockRepo.On("GetAll").Return(expectedUsers, nil)

	users, err := service.GetUsers()

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers[0].Name, users[0].Name)
	assert.Equal(t, expectedUsers[1].Name, users[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		request        *models.UpdateUserRequest
		mockSetup      func(*MockUserRepository)
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:   "Successfully update user",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			request: &models.UpdateUserRequest{
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				existingUser := &models.User{
					ID:    userID,
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}
				mockRepo.On("GetByID", userID).Return(existingUser, nil)
				mockRepo.On("GetByEmail", "john.updated@example.com").Return(nil, fmt.Errorf("user not found"))
				mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "Email already taken by another user",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			request: &models.UpdateUserRequest{
				Name:  "John Updated",
				Email: "jane.smith@example.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				existingUser := &models.User{
					ID:    userID,
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}
				otherUser := &models.User{
					ID:    uuid.New(),
					Name:  "Jane Smith",
					Email: "jane.smith@example.com",
				}
				mockRepo.On("GetByID", userID).Return(existingUser, nil)
				mockRepo.On("GetByEmail", "jane.smith@example.com").Return(otherUser, nil)
			},
			expectedError:  true,
			expectedErrMsg: "user with email jane.smith@example.com already exists",
		},
		{
			name:   "User not found",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			request: &models.UpdateUserRequest{
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			mockSetup: func(mockRepo *MockUserRepository) {
				userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				mockRepo.On("GetByID", userID).Return(nil, fmt.Errorf("user not found"))
			},
			expectedError:  true,
			expectedErrMsg: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUserService(mockRepo)

			tt.mockSetup(mockRepo)

			user, err := service.UpdateUser(tt.userID, tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.request.Name, user.Name)
				assert.Equal(t, tt.request.Email, user.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*MockUserRepository)
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:   "Successfully delete user",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(mockRepo *MockUserRepository) {
				userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				existingUser := &models.User{
					ID:    userID,
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}
				mockRepo.On("GetByID", userID).Return(existingUser, nil)
				mockRepo.On("Delete", userID).Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "User not found",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(mockRepo *MockUserRepository) {
				userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				mockRepo.On("GetByID", userID).Return(nil, fmt.Errorf("user not found"))
			},
			expectedError:  true,
			expectedErrMsg: "user not found",
		},
		{
			name:   "Invalid UUID format",
			userID: "invalid-uuid",
			mockSetup: func(mockRepo *MockUserRepository) {
				// No mock setup needed for invalid UUID
			},
			expectedError:  true,
			expectedErrMsg: "invalid user ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := NewUserService(mockRepo)

			tt.mockSetup(mockRepo)

			err := service.DeleteUser(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
