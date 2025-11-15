package e2e

import (
	"database/sql"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	accountapp "github.com/raihanstark/trade-journal/internal/application/account"
	"github.com/raihanstark/trade-journal/internal/application/auth"
	strategyapp "github.com/raihanstark/trade-journal/internal/application/strategy"
	tradeapp "github.com/raihanstark/trade-journal/internal/application/trade"
	"github.com/raihanstark/trade-journal/internal/db"
	"github.com/raihanstark/trade-journal/internal/infrastructure/http/handlers"
	custommiddleware "github.com/raihanstark/trade-journal/internal/infrastructure/http/middleware"
	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
	"github.com/raihanstark/trade-journal/internal/infrastructure/security"
)

// setupTestServer creates a fully configured Echo server for e2e tests
func setupTestServer(t *testing.T, _ *sql.DB, queries *db.Queries) *echo.Echo {
	t.Helper()

	// Initialize infrastructure
	userRepository := persistence.NewUserRepository(queries)
	accountRepository := persistence.NewAccountRepository(queries)
	strategyRepository := persistence.NewStrategyRepository(queries)
	tradeRepository := persistence.NewTradeRepository(queries)
	tokenGenerator := security.NewJWTTokenGenerator("test-secret-key")

	// Initialize application layer
	authService := auth.NewService(userRepository, tokenGenerator)
	accountService := accountapp.NewService(accountRepository)
	strategyService := strategyapp.NewService(strategyRepository)
	tradeService := tradeapp.NewService(tradeRepository, accountRepository)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	accountHandler := handlers.NewAccountHandler(accountService)
	strategyHandler := handlers.NewStrategyHandler(strategyService)
	tradeHandler := handlers.NewTradeHandler(tradeService)

	// Create Echo instance
	e := echo.New()
	e.Use(middleware.Recover())

	// Public routes
	e.POST("/api/auth/register", authHandler.Register)
	e.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	protected := e.Group("/api")
	protected.Use(custommiddleware.JWTAuth(tokenGenerator))

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

	// Trade routes
	protected.POST("/trades", tradeHandler.CreateTrade)
	protected.GET("/trades", tradeHandler.GetTrades)
	protected.GET("/trades/:id", tradeHandler.GetTrade)
	protected.PUT("/trades/:id", tradeHandler.UpdateTrade)
	protected.DELETE("/trades/:id", tradeHandler.DeleteTrade)

	return e
}
