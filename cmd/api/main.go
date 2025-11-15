package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	accountapp "github.com/raihanstark/trade-journal/internal/application/account"
	"github.com/raihanstark/trade-journal/internal/application/auth"
	strategyapp "github.com/raihanstark/trade-journal/internal/application/strategy"
	"github.com/raihanstark/trade-journal/internal/db"
	"github.com/raihanstark/trade-journal/internal/infrastructure/http/handlers"
	custommiddleware "github.com/raihanstark/trade-journal/internal/infrastructure/http/middleware"
	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
	"github.com/raihanstark/trade-journal/internal/infrastructure/security"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get configuration from environment
	databaseURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	port := os.Getenv("PORT")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	if port == "" {
		port = "8080"
	}

	// Connect to database
	dbConn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Test database connection
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Initialize infrastructure layer
	queries := db.New(dbConn)
	userRepository := persistence.NewUserRepository(queries)
	accountRepository := persistence.NewAccountRepository(queries)
	strategyRepository := persistence.NewStrategyRepository(queries)
	tokenGenerator := security.NewJWTTokenGenerator(jwtSecret)

	// Initialize application layer
	authService := auth.NewService(userRepository, tokenGenerator)
	accountService := accountapp.NewService(accountRepository)
	strategyService := strategyapp.NewService(strategyRepository)

	// Initialize presentation layer
	authHandler := handlers.NewAuthHandler(authService)
	accountHandler := handlers.NewAccountHandler(accountService)
	strategyHandler := handlers.NewStrategyHandler(strategyService)

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Public routes
	e.POST("/api/auth/register", authHandler.Register)
	e.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	protected := e.Group("/api")
	protected.Use(custommiddleware.JWTAuth(tokenGenerator))

	// User routes
	protected.GET("/me", func(c echo.Context) error {
		userID := c.Get("user_id").(int64)
		email := c.Get("email").(string)
		return c.JSON(200, map[string]any{
			"user_id": userID,
			"email":   email,
		})
	})

	// Account routes
	protected.POST("/accounts", accountHandler.CreateAccount)
	protected.GET("/accounts", accountHandler.GetAccounts)
	protected.GET("/accounts/:id", accountHandler.GetAccount)
	protected.PUT("/accounts/:id", accountHandler.UpdateAccount)
	protected.DELETE("/accounts/:id", accountHandler.DeleteAccount)

	// Strategy routes
	protected.POST("/strategies", strategyHandler.CreateStrategy)
	protected.GET("/strategies", strategyHandler.GetStrategies)
	protected.GET("/strategies/:id", strategyHandler.GetStrategy)
	protected.PUT("/strategies/:id", strategyHandler.UpdateStrategy)
	protected.DELETE("/strategies/:id", strategyHandler.DeleteStrategy)

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
