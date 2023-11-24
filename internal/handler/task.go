package handler

import (
	"net/http"

	"github.com/Irori235/system-design-2023-v2/internal/repository"
	"github.com/gin-gonic/gin"
	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type (
	GetTasksResponse []GetTaskResponse
	GetTaskResponse  struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		Title     string    `json:"title"`
		IsDone    bool      `json:"is_done"`
		CreatedAt string    `json:"created_at"`
	}

	SearchTasksRequest struct {
	}

	CreateTaskRequest struct {
		Title string `json:"title"`
	}

	UpdateTaskRequest struct {
		Title  string `json:"title"`
		IsDone bool   `json:"is_done"`
	}
)

// GET /api/v1/tasks
func (h *Handler) GetTasks(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	tasks, err := h.repo.GetTasks(c, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make(GetTasksResponse, len(tasks))
	for i, task := range tasks {
		res[i] = GetTaskResponse{
			ID:        task.ID,
			UserID:    task.UserID,
			Title:     task.Title,
			IsDone:    task.IsDone,
			CreatedAt: task.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, res)

}

// POST /api/v1/tasks
func (h *Handler) CreateTask(c *gin.Context) {
	req := new(CreateTaskRequest)
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	params := repository.CreateTaskParams{
		UserID: userID.(uuid.UUID),
		Title:  req.Title,
	}

	err := h.repo.CreateTask(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// PUT /api/v1/tasks/:taskID
func (h *Handler) UpdateTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := new(UpdateTaskRequest)
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = vd.ValidateStruct(
		req,
		vd.Field(&req.Title, vd.Required),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := repository.UpdateTaskParams{
		ID:     taskID,
		Title:  req.Title,
		IsDone: req.IsDone,
	}

	err = h.repo.UpdateTask(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DELETE /api/v1/tasks/:taskID
func (h *Handler) DeleteTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.repo.DeleteTask(c, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
