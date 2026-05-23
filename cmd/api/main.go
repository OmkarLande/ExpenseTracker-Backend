package main

import (
	"fmt"
	"log"
	"net/http"

	"ExpenseTracker-Backend/internal/config"
	"ExpenseTracker-Backend/internal/database"
	domain_auth "ExpenseTracker-Backend/internal/domain/auth"
	domain_transaction "ExpenseTracker-Backend/internal/domain/transaction"
	domain_user "ExpenseTracker-Backend/internal/domain/user"
	"ExpenseTracker-Backend/internal/middleware"
	"ExpenseTracker-Backend/internal/storage"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Connect to database
	dbpool, err := database.NewPostgresConnection(cfg.DBUrl)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer dbpool.Close()

	// Initialize S3 Storage Wrapper
	s3Storage, err := storage.NewS3Storage(cfg)
	if err != nil {
		log.Printf("[WARNING] S3 Storage initialization failed: %v", err)
	}

	// Initialize components
	emailSender := domain_auth.NewEmailSender(cfg)
	authRepo := domain_auth.NewRepository(dbpool)
	authService := domain_auth.NewService(authRepo, emailSender, cfg)
	authHandler := domain_auth.NewHandler(authService)

	// Initialize User Domain
	userRepo := domain_user.NewRepository(dbpool)
	userService := domain_user.NewService(userRepo, s3Storage, cfg)
	userHandler := domain_user.NewHandler(userService)

	// Initialize Transaction Domain
	txRepo := domain_transaction.NewRepository(dbpool)
	txService := domain_transaction.NewService(txRepo)
	txHandler := domain_transaction.NewHandler(txService)

	// Setup Chi Router
	r := chi.NewRouter()

	// Apply global middlewares
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	// Register routes
	r.Route("/api/v1", func(r chi.Router) {
		domain_auth.RegisterRoutes(r, authHandler, authMiddleware)
		domain_user.RegisterRoutes(r, userHandler, authMiddleware)
		domain_transaction.RegisterRoutes(r, txHandler, authMiddleware)
	})

	// Start server
	fmt.Println("Server running on port:", cfg.Port)
	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Server error:", err)
	}
}