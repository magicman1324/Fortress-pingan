package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/pingan/bastion/internal/api/ws"
	"github.com/pingan/bastion/internal/service"
)

func RegisterRoutes(
	router *gin.Engine,
	authSvc *service.AuthService,
	assetSvc *service.AssetService,
	sessionSvc *service.SessionService,
	auditSvc *service.AuditService,
	wsHandler *ws.Handler,
	wsHub *ws.Hub,
	jwtSecret string,
) {
	authH := NewAuthHandler(authSvc)
	assetH := NewAssetHandler(assetSvc)
	sessionH := NewSessionHandler(sessionSvc)
	auditH := NewAuditHandler(auditSvc)
	dashboardH := NewDashboardHandler(sessionSvc, auditSvc)

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", authH.Login)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Authenticated routes
	authed := router.Group("/api/v1")
	authed.Use(AuthMiddleware(jwtSecret))
	authed.Use(TokenValidationMiddleware(authSvc))
	{
		authed.GET("/auth/me", authH.Me)

		authed.GET("/assets", assetH.List)
		authed.POST("/assets", assetH.Create)
		authed.PUT("/assets/:id", assetH.Update)
		authed.DELETE("/assets/:id", assetH.Delete)
		authed.POST("/assets/:id/test", assetH.TestConnect)

		authed.GET("/sessions", sessionH.List)
		authed.POST("/sessions/:id/terminate", sessionH.Terminate)

		authed.GET("/audit-logs", auditH.Search)

		authed.GET("/dashboard/stats", dashboardH.Stats)
	}

	// WebSocket terminal (JWT via query param for browser WebSocket)
	router.GET("/api/v1/ws/terminal/:assetId", WSJWTMiddleware(authSvc), func(c *gin.Context) {
		wsHandler.HandleTerminal(c, wsHub)
	})
}
