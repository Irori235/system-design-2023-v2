package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Irori235/system-design-2023-v2/internal/repository"
	"github.com/gin-gonic/gin"
	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type (
	SignUpRequest struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	SignInRequest struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	Claims struct {
		UserID string `json:"user_id"`
		jwt.StandardClaims
	}

	SignUpResponse struct {
		ID uuid.UUID `json:"id"`
	}

	SignInResponse struct {
		Token string `json:"token"`
	}
)

func (h *Handler) SignUp(c *gin.Context) {
	req := new(SignUpRequest)
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := vd.ValidateStruct(
		req,
		vd.Field(&req.Name, vd.Required),
		vd.Field(&req.Password, vd.Required),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request body: %w", err).Error()})
		return
	}

	params := repository.CreateUserParams{
		Name:     req.Name,
		Password: req.Password,
	}

	userID, err := h.repo.CreateUser(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := SignUpResponse{
		ID: userID,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) SignIn(c *gin.Context) {
	req := new(SignInRequest)
	if err := c.Bind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := vd.ValidateStruct(
		req,
		vd.Field(&req.Name, vd.Required),
		vd.Field(&req.Password, vd.Required),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request body: %w", err).Error()})
		return
	}

	userID, err := h.repo.GetUserID(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ok, err := h.repo.CheckPass(c, userID, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id or password"})
		return
	}

	token, err := generateJWT(userID.String(), h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// res := SignInResponse{
	// 	Token: token,
	// }

	// set cookie
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(3 * time.Hour),
		HttpOnly: true,
		// Secure:   true,
		Path: "/",
		// SameSite: http.SameSiteStrictMode,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, cookie)
	// http.Header.Add(c.Writer.Header(), "Access-Control-Allow-Credentials", "true")

	// c.JSON(http.StatusOK, res)
}

func (h *Handler) SignOut(c *gin.Context) {
	// delete cookie
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "expired",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		// Secure:   true,
		Path: "/",
		// SameSite: http.SameSiteStrictMode,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, cookie)
	// http.Header.Add(c.Writer.Header(), "Access-Control-Allow-Credentials", "true")

	c.JSON(http.StatusOK, gin.H{})
}

func generateJWT(userID string, jwtSecret string) (string, error) {

	expiresAt := time.Now().Add(3 * time.Hour)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("generate jwt: %w", err)
	}

	return tokenStr, nil

}
