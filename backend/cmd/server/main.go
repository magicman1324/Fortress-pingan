package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/pingan/bastion/internal/api/rest"
	"github.com/pingan/bastion/internal/api/ws"
	"github.com/pingan/bastion/internal/repository"
	"github.com/pingan/bastion/internal/service"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	dsn := getEnv("DB_DSN", "root:@tcp(127.0.0.1:4000)/bastion?parseTime=true")
	encryptionKey := getEnv("ENCRYPTION_KEY", "bastion-dev-key-32bytes-long!!")
	jwtSecret := getEnv("JWT_SECRET", "bastion-jwt-secret-key")

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Repositories
	userRepo := repository.NewUserRepo(db)
	assetRepo := repository.NewAssetRepo(db)
	sessionRepo := repository.NewSessionRepo(db)
	auditRepo := repository.NewAuditRepo(db)

	// Services
	authSvc := service.NewAuthService(userRepo, jwtSecret)
	assetSvc := service.NewAssetService(assetRepo, encryptionKey)
	sessionSvc := service.NewSessionService(sessionRepo)
	auditSvc := service.NewAuditService(auditRepo)

	// WebSocket handler
	wsHandler := ws.NewHandler(assetSvc, sessionSvc, auditRepo, jwtSecret)
	wsHub := ws.NewHub()
	go wsHub.Run()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	rest.RegisterRoutes(router, authSvc, assetSvc, sessionSvc, auditSvc, wsHandler, wsHub, jwtSecret)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:   60 * time.Second,
	}

	go func() {
		log.Printf("Bastion backend listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}
	log.Println("Server stopped")
}
