package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manishrw/sample-crud-api/go-echo/internal/models"
	"github.com/manishrw/sample-crud-api/go-echo/internal/services"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers handles GET /api/v1/users
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.userService.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// GetUser handles GET /api/v1/users/:id
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "User ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "User not found",
				Code:    http.StatusNotFound,
			})
		}
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
	}

	// Basic validation
	if req.Name == "" || req.Email == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Name and email are required",
			Code:    http.StatusBadRequest,
		})
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		if err.Error() == "user with email "+req.Email+" already exists" {
			return c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "Conflict",
				Message: err.Error(),
				Code:    http.StatusConflict,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "User ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request body",
			Code:    http.StatusBadRequest,
		})
	}

	// Basic validation
	if req.Name == "" || req.Email == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Name and email are required",
			Code:    http.StatusBadRequest,
		})
	}

	user, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "User not found",
				Code:    http.StatusNotFound,
			})
		}
		if err.Error() == "user with email "+req.Email+" already exists" {
			return c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "Conflict",
				Message: err.Error(),
				Code:    http.StatusConflict,
			})
		}
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "User ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	err := h.userService.DeleteUser(id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Not Found",
				Message: "User not found",
				Code:    http.StatusNotFound,
			})
		}
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User deleted successfully",
	})
}
