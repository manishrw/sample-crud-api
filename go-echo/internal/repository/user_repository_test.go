package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/manishrw/sample-crud-api/go-echo/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *UserRepository) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	repo := NewUserRepository(db)
	return db, mock, repo
}

func TestCreateUser(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	user := &models.User{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	// Expect the INSERT query
	mock.ExpectExec("INSERT INTO users").
		WithArgs(sqlmock.AnyArg(), user.Name, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(user)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	userID := uuid.New()
	expectedUser := &models.User{
		ID:        userID,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetByID(userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	userID := uuid.New()

	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID(userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	expectedUsers := []*models.User{
		{
			ID:        uuid.New(),
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Jane Smith",
			Email:     "jane.smith@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"})
	for _, user := range expectedUsers {
		rows.AddRow(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	}

	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC").
		WillReturnRows(rows)

	users, err := repo.GetAll()

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers[0].Name, users[0].Name)
	assert.Equal(t, expectedUsers[1].Name, users[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	user := &models.User{
		ID:    uuid.New(),
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	mock.ExpectExec("UPDATE users SET name = \\$1, email = \\$2, updated_at = \\$3 WHERE id = \\$4").
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(user)

	assert.NoError(t, err)
	assert.False(t, user.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_NotFound(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	user := &models.User{
		ID:    uuid.New(),
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}

	mock.ExpectExec("UPDATE users SET name = \\$1, email = \\$2, updated_at = \\$3 WHERE id = \\$4").
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	userID := uuid.New()

	mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser_NotFound(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	userID := uuid.New()

	mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Delete(userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmail(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	email := "john.doe@example.com"
	expectedUser := &models.User{
		ID:        uuid.New(),
		Name:      "John Doe",
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmail_NotFound(t *testing.T) {
	db, mock, repo := setupTestDB(t)
	defer db.Close()

	email := "nonexistent@example.com"

	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users WHERE email = \\$1").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByEmail(email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}
