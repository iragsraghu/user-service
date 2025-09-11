package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/iragsraghu/user-service/api/v1"
	"github.com/iragsraghu/user-service/internal/config"
	"github.com/iragsraghu/user-service/internal/db"
	"github.com/iragsraghu/user-service/internal/handler"
	"github.com/iragsraghu/user-service/internal/models"
	"github.com/iragsraghu/user-service/internal/repo"
)

func main() {
	cfg := config.Load()

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	dbConn, err := db.NewMySQL(user, pass, host, port, name)
	fmt.Println(dbConn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// AutoMigrate schema
	if err := dbConn.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	userRepo := repo.NewUserRepo(dbConn)
	userHandler := handler.NewUserHandler(userRepo)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Middleware: context timeout
	r.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	v1.Register(r, userHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("server stopped")
}
