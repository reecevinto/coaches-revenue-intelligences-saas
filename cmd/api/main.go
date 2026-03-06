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
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/auth"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/config"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/database"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/jwt"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/middleware"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/users"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/pkg/logger"
)

func main() {

	// =============================
	// 1️⃣ Load Configuration
	// =============================

	cfg := config.Load()

	// initialize structured logger
	appLogger := logger.New(cfg.AppEnv)
	log.Logger = appLogger

	// =============================
	// 2️⃣ Initialize Infrastructure
	// =============================

	// create postgres connection pool
	dbPool, err := database.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer dbPool.Close()

	// load RSA keys for JWT signing
	jwtService, err := jwt.NewService("keys/private.pem", "keys/public.pem")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load jwt keys")
	}

	// =============================
	// 3️⃣ Initialize Repositories
	// =============================

	accountRepo := database.NewAccountRepository(dbPool)
	userRepo := database.NewUserRepository(dbPool)

	// =============================
	// 4️⃣ Initialize Services
	// =============================

	accountService := accounts.NewService(accountRepo)
	userService := users.NewService(userRepo)
	authService := auth.NewService(userRepo, jwtService)

	// prevent unused warnings until handlers use them
	_ = accountService

	// =============================
	// 5️⃣ Setup Router
	// =============================

	r := chi.NewRouter()

	// global middleware
	r.Use(middleware.RequestID)

	// health check
	r.Get("/health", middleware.HealthHandler)

	// -------- USERS ROUTES --------
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(r)

	// -------- AUTH ROUTES --------
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRoutes(r)

	// =============================
	// 6️⃣ HTTP Server Configuration
	// =============================

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// run server in goroutine
	go func() {

		log.Info().
			Str("port", cfg.Port).
			Msg("server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().
				Err(err).
				Msg("server failed")
		}

	}()

	// =============================
	// 7️⃣ Graceful Shutdown
	// =============================

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Info().Msg("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("server forced shutdown")
	}

	log.Info().Msg("server exited cleanly")

}
