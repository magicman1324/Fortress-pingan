package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pingan/bastion/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}
	token, user, err := h.authSvc.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	token := c.GetString("token")
	userID, role, err := h.authSvc.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"role":    role,
	})
}
