package handler

import (
	"crypto/rand"

	"github.com/Irori235/system-design-2023-v2/internal/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	jwtSecret string
	repo      *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	jwtSecret := randomString()

	return &Handler{
		jwtSecret: jwtSecret,
		repo:      repo,
	}
}

func (h *Handler) SetupRoutes(group *gin.RouterGroup) {
	// ping group
	pingAPI := group.Group("/ping")
	{
		pingAPI.GET("", h.Ping)
	}

	// user group
	userAPI := group.Group("/users")
	userAPI.Use(h.AuthMiddleware())
	{
		userAPI.GET("/me", h.GetMe)
		userAPI.PATCH("/name", h.UpdateName)
		userAPI.PATCH("/password", h.UpdatePass)
		userAPI.DELETE("/quit", h.Quit)
	}

	// task group
	taskAPI := group.Group("/tasks")
	taskAPI.Use(h.AuthMiddleware())
	{
		taskAPI.GET("", h.GetTasks)
		taskAPI.GET("/search", h.SearchTasks)
		taskAPI.POST("", h.CreateTask)
		taskAPI.PUT("/:taskID", h.UpdateTask)
		taskAPI.DELETE("/:taskID", h.DeleteTask)
	}

	// auth group
	authAPI := group.Group("/auth")
	{
		authAPI.POST("/signup", h.SignUp)
		authAPI.POST("/signin", h.SignIn)
	}
}

func randomString() string {
	length := 32
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	for i := 0; i < length; i++ {
		b[i] = letters[b[i]%byte(len(letters))]
	}

	return string(b)
}
