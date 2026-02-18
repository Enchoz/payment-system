package main

import (
	"fmt"
	"log"
	"os"
	"payment-system/internal/application/services"
	"payment-system/internal/infrastructure/postgres"
	"payment-system/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		dbConnStr = "postgres://postgres:postgres@localhost:5432/payment_system?sslmode=disable"
	}

	db, err := postgres.NewDB(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	externalAccountRepo := postgres.NewExternalAccountRepository(db)
	paymentRequestRepo := postgres.NewPaymentRequestRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	ledgerRepo := postgres.NewLedgerRepository(db)
	paymentHoldRepo := postgres.NewPaymentHoldRepository(db)

	// Initialize services
	exchangeRateService := services.NewExchangeRateService(ledgerRepo)
	validationService := services.NewValidationService()
	paymentService := services.NewPaymentService(
		db,
		userRepo,
		accountRepo,
		externalAccountRepo,
		paymentRequestRepo,
		paymentRepo,
		ledgerRepo,
		paymentHoldRepo,
		exchangeRateService,
		validationService,
	)

	// Initialize handlers
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Setup router
	router := gin.Default()

	// API routes
	v1 := router.Group("/api/v1")
	{
		payments := v1.Group("/payments")
		{
			payments.POST("/internal", paymentHandler.ProcessInternalPayment)
			payments.POST("/external", paymentHandler.CreateExternalPaymentRequest)
			payments.POST("/requests/:id/process", paymentHandler.ProcessPaymentRequest)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "4078"
	}

	fmt.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
