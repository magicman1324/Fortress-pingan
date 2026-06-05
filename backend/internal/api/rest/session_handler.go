package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingan/bastion/internal/model"
	"github.com/pingan/bastion/internal/service"
)

type SessionHandler struct {
	svc *service.SessionService
}

func NewSessionHandler(svc *service.SessionService) *SessionHandler {
	return &SessionHandler{svc: svc}
}

func (h *SessionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	offset := (page - 1) * size
	sessions, err := h.svc.List(size, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sessions == nil {
		sessions = []model.Session{}
	}
	c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) Terminate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	// Termination is handled by the WS hub; here we just close in DB
	if err := h.svc.Close(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "terminated"})
}
