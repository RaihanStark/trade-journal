package e2e

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	accountapp "github.com/raihanstark/trade-journal/internal/application/account"
	analyticsapp "github.com/raihanstark/trade-journal/internal/application/analytics"
	"github.com/raihanstark/trade-journal/internal/application/auth"
	strategyapp "github.com/raihanstark/trade-journal/internal/application/strategy"
	tradeapp "github.com/raihanstark/trade-journal/internal/application/trade"
	"github.com/raihanstark/trade-journal/internal/db"
	"github.com/raihanstark/trade-journal/internal/infrastructure/http/handlers"
	custommiddleware "github.com/raihanstark/trade-journal/internal/infrastructure/http/middleware"
	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
	"github.com/raihanstark/trade-journal/internal/infrastructure/security"
	"github.com/raihanstark/trade-journal/internal/infrastructure/storage"
)

// setupTestServer creates a fully configured Echo server for e2e tests
func setupTestServer(t *testing.T, _ *sql.DB, queries *db.Queries) *echo.Echo {
	t.Helper()

	// Initialize infrastructure
	userRepository := persistence.NewUserRepository(queries)
	accountRepository := persistence.NewAccountRepository(queries)
	strategyRepository := persistence.NewStrategyRepository(queries)
	tradeRepository := persistence.NewTradeRepository(queries)
	analyticsRepository := persistence.NewAnalyticsRepository(queries)
	tokenGenerator := security.NewJWTTokenGenerator("test-secret-key")

	// Initialize application layer
	authService := auth.NewService(userRepository, tokenGenerator)
	accountService := accountapp.NewService(accountRepository)
	strategyService := strategyapp.NewService(strategyRepository)
	tradeService := tradeapp.NewService(tradeRepository, accountRepository)
	analyticsService := analyticsapp.NewService(analyticsRepository)

	// Initialize storage (MinIO for tests)
	minioStorage, err := storage.NewMinIOStorage("localhost:9000", "minioadmin", "minioadmin123", "trade-journal", false)
	if err != nil {
		t.Fatalf("failed to create MinIO storage: %v", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	err = minioStorage.EnsureBucket(ctx)
	if err != nil {
		t.Fatalf("failed to ensure MinIO bucket exists: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	accountHandler := handlers.NewAccountHandler(accountService)
	strategyHandler := handlers.NewStrategyHandler(strategyService)
	tradeHandler := handlers.NewTradeHandler(tradeService, minioStorage)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

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
	protected.POST("/trades/:id/chart/:type", tradeHandler.UploadChart)

	// Analytics routes
	protected.GET("/analytics", analyticsHandler.GetUserAnalytics)

	return e
}

// createTestImage creates a simple 1x1 PNG image for testing
func createTestImage(t *testing.T) []byte {
	t.Helper()

	// Create a 1x1 image
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255}) // Red pixel

	// Encode to PNG
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("failed to encode test image: %v", err)
	}

	return buf.Bytes()
}

// createMultipartForm creates a multipart form with a file upload
func createMultipartForm(t *testing.T, body io.Writer, fieldName, filename string, fileData []byte) *multipart.Writer {
	t.Helper()

	writer := multipart.NewWriter(body)

	// Detect content type from filename
	contentType := "image/png"
	if len(filename) > 4 && filename[len(filename)-4:] == ".txt" {
		contentType = "text/plain"
	}

	// Create form file with proper content type header
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filename))
	h.Set("Content-Type", contentType)

	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	_, err = part.Write(fileData)
	if err != nil {
		t.Fatalf("failed to write file data: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	return writer
}
