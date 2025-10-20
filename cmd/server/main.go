package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"emplopyee-app-go/internal/config"
	"emplopyee-app-go/internal/dao"
	"emplopyee-app-go/internal/db"
	"emplopyee-app-go/internal/router"
	"emplopyee-app-go/internal/service"
)

func main() {
	cfg := config.Load() // reads from env/defaults

	// Initialize DB pool
	pool, err := db.NewDB(cfg.DatabaseDSN, cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}
	defer pool.Close()

	// Wire dependencies (manual DI)
	empDAO := dao.NewEmployeeDAO(pool)
	empService := service.NewEmployeeService(empDAO)
	r := router.NewRouter(empService)

	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s\n", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server Shutdown: %v", err)
	}
	log.Println("server stopped")
}
