package handler

import (
	"fmt"
	"net/http"

	"github.com/Irori235/system-design-2023-v2/internal/repository"

	"github.com/gin-gonic/gin"
	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

// スキーマ定義
type (
	GetUsersResponse []GetUserResponse

	GetUserResponse struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}

	CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	CreateUserResponse struct {
		ID uuid.UUID `json:"id"`
	}
)

// GET /api/v1/users
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.repo.GetUsers(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make(GetUsersResponse, len(users))
	for i, user := range users {
		res[i] = GetUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	c.JSON(http.StatusOK, res)
}

// POST /api/v1/users
func (h *Handler) CreateUser(c *gin.Context) {
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := vd.ValidateStruct(
		req,
		vd.Field(&req.Name, vd.Required),
		vd.Field(&req.Email, vd.Required, is.Email),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request body: %w", err).Error()})
		return
	}

	params := repository.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
	}

	userID, err := h.repo.CreateUser(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := CreateUserResponse{
		ID: userID,
	}

	c.JSON(http.StatusOK, res)
}

// GET /api/v1/users/:userID
func (h *Handler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid userID: %w", err).Error()})
		return
	}

	user, err := h.repo.GetUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := GetUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusOK, res)
}
