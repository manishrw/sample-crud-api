package services

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/manishrw/sample-crud-api/go-echo/internal/models"
	"github.com/manishrw/sample-crud-api/go-echo/internal/repository"
)

// UserServiceInterface defines the contract for user service operations
type UserServiceInterface interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUser(id string) (*models.User, error)
	GetUsers() ([]*models.User, error)
	UpdateUser(id string, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id string) error
}

type UserService struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Validate email uniqueness
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		Name:  strings.TrimSpace(req.Name),
		Email: strings.ToLower(strings.TrimSpace(req.Email)),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUsers() ([]*models.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (s *UserService) UpdateUser(id string, req *models.UpdateUserRequest) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Check if user exists
	existingUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if email is being changed and if it's already taken
	if existingUser.Email != strings.ToLower(strings.TrimSpace(req.Email)) {
		userWithEmail, err := s.userRepo.GetByEmail(req.Email)
		if err == nil && userWithEmail != nil && userWithEmail.ID != userID {
			return nil, fmt.Errorf("user with email %s already exists", req.Email)
		}
	}

	// Update user fields
	existingUser.Name = strings.TrimSpace(req.Name)
	existingUser.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := s.userRepo.Update(existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return existingUser, nil
}

func (s *UserService) DeleteUser(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Check if user exists
	_, err = s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := s.userRepo.Delete(userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
