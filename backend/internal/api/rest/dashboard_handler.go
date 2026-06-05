package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pingan/bastion/internal/service"
)

type DashboardHandler struct {
	sessionSvc *service.SessionService
	auditSvc   *service.AuditService
}

func NewDashboardHandler(sessionSvc *service.SessionService, auditSvc *service.AuditService) *DashboardHandler {
	return &DashboardHandler{sessionSvc: sessionSvc, auditSvc: auditSvc}
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	activeSessions, _ := h.sessionSvc.ActiveSessionCount()
	totalSessions, _ := h.sessionSvc.TotalSessionCount()
	todayAuditCount, _ := h.getTodayAuditCount()
	recentAudits, _, _ := h.auditSvc.Search(0, 0, "", 1, 5)

	c.JSON(http.StatusOK, gin.H{
		"active_sessions":  activeSessions,
		"total_sessions":   totalSessions,
		"today_audit_count": todayAuditCount,
		"recent_audits":     recentAudits,
	})
}

func (h *DashboardHandler) getTodayAuditCount() (int, error) {
	_, total, err := h.auditSvc.Search(0, 0, "", 1, 1)
	return total, err
}
