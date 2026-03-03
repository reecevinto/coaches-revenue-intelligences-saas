package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/accounts"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/config"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/database"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/middleware"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/users"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/pkg/logger"
)

func main() {

	// =============================
	// 1️⃣  Load Configuration
	// =============================
	cfg := config.Load()

	appLogger := logger.New(cfg.AppEnv)
	log.Logger = appLogger

	// =============================
	// 2️⃣  Initialize Infrastructure
	// =============================
	dbPool, err := database.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer dbPool.Close()

	// =============================
	// 3️⃣  Initialize Repositories
	// =============================
	accountRepo := database.NewAccountRepository(dbPool)
	userRepo := database.NewUserRepository(dbPool)

	// =============================
	// 4️⃣  Initialize Services
	// =============================
	accountService := accounts.NewService(accountRepo)
	userService := users.NewService(userRepo)

	// Temporary until handlers are built
	_ = accountService
	_ = userService

	// =============================
	// 5️⃣  Setup Router
	// =============================
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	r.Get("/health", middleware.HealthHandler)

	// ===== Add Users Handler =====
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(r)

	// =============================
	// 6️⃣  Start HTTP Server
	// =============================
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Info().Msg("server starting on port " + cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	// =============================
	// 7️⃣  Graceful Shutdown
	// =============================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exited properly")
}
