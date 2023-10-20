package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Irori235/system-design-2023-v2/internal/repository"

	"github.com/gin-gonic/gin"
	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

// スキーマ定義
type (
	UpdateNameRequest struct {
		Name string `json:"name"`
	}

	UpdatePassRequest struct {
		Password string `json:"password"`
	}

	GetMeResponse struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
	}
)

// GET /api/v1/users/me
func (h *Handler) GetMe(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	user, err := h.repo.GetUser(c, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := GetMeResponse{
		ID:        user.ID,
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

// PATCH /api/v1/users/name
func (h *Handler) UpdateName(c *gin.Context) {
	req := new(UpdateNameRequest)
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := vd.ValidateStruct(
		req,
		vd.Field(&req.Name, vd.Required),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request body: %w", err).Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	params := repository.UpdateNameParams{
		ID:   userID.(uuid.UUID),
		Name: req.Name,
	}

	err = h.repo.UpdateName(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// PATCH /api/v1/users/password
func (h *Handler) UpdatePass(c *gin.Context) {
	req := new(UpdatePassRequest)
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := vd.ValidateStruct(
		req,
		vd.Field(&req.Password, vd.Required),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request body: %w", err).Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	params := repository.UpdatePassParams{
		ID:       userID.(uuid.UUID),
		Password: req.Password,
	}

	err = h.repo.UpdatePass(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DELETE /api/v1/users/quit
func (h *Handler) Quit(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	err := h.repo.DeleteUser(c, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
