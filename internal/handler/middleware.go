package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieHeader := c.GetHeader("cookie")
		if cookieHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cookie header is required"})
			c.Abort()
			return
		}

		cookies := strings.Split(cookieHeader, ";")
		var tokenString string
		for _, cookie := range cookies {
			if strings.Contains(cookie, "jwt=") {
				tokenString = strings.Replace(cookie, "jwt=", "", 1)
			}
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "jwt cookie is required"})
			c.Abort()
			return
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.jwtSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		str := claims.UserID
		userID, err := uuid.Parse(str)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
