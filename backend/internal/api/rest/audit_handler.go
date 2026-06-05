package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingan/bastion/internal/model"
	"github.com/pingan/bastion/internal/service"
)

type AuditHandler struct {
	svc *service.AuditService
}

func NewAuditHandler(svc *service.AuditService) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (h *AuditHandler) Search(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	userID, _ := strconv.ParseInt(c.DefaultQuery("user_id", "0"), 10, 64)
	assetID, _ := strconv.ParseInt(c.DefaultQuery("asset_id", "0"), 10, 64)
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	logs, total, err := h.svc.Search(userID, assetID, keyword, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if logs == nil {
		logs = []model.AuditLog{}
	}
	c.JSON(http.StatusOK, gin.H{
		"items": logs,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
