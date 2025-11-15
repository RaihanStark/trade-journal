package strategy

import (
	"context"
	"testing"

	"github.com/raihanstark/trade-journal/internal/domain/user"
	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
	"github.com/raihanstark/trade-journal/internal/testutil"
)

// Integration tests for strategy service
// Even though it's simple CRUD, we verify database operations work correctly

func TestStrategyService_CreateStrategy_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	strategyRepo := persistence.NewStrategyRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(strategyRepo)

	ctx := context.Background()

	t.Run("creates strategy in database", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("strategy@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		// Create strategy
		req := CreateStrategyRequest{
			Name:        "Breakout Strategy",
			Description: "Trade price breakouts at key levels",
		}

		result, err := service.CreateStrategy(ctx, createdUser.ID, req)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify response
		if result.Name != "Breakout Strategy" {
			t.Errorf("expected name 'Breakout Strategy', got %s", result.Name)
		}
		if result.Description != "Trade price breakouts at key levels" {
			t.Errorf("expected description 'Trade price breakouts at key levels', got %s", result.Description)
		}

		// Verify strategy was saved to database
		var savedName, savedDesc string
		err = pg.DB.QueryRow("SELECT name, description FROM strategies WHERE id = $1", result.ID).
			Scan(&savedName, &savedDesc)

		if err != nil {
			t.Fatalf("failed to query saved strategy: %v", err)
		}

		if savedName != "Breakout Strategy" {
			t.Errorf("expected saved name 'Breakout Strategy', got %s", savedName)
		}
		if savedDesc != "Trade price breakouts at key levels" {
			t.Errorf("expected saved description, got %s", savedDesc)
		}
	})
}

func TestStrategyService_GetStrategy_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	strategyRepo := persistence.NewStrategyRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(strategyRepo)

	ctx := context.Background()

	t.Run("retrieves strategy from database", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("getstrategy@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create strategy
		req := CreateStrategyRequest{
			Name:        "Support/Resistance",
			Description: "Trade at S/R levels",
		}
		created, _ := service.CreateStrategy(ctx, createdUser.ID, req)

		// Retrieve strategy
		retrieved, err := service.GetStrategy(ctx, created.ID, createdUser.ID)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify retrieved data matches created data
		if retrieved.ID != created.ID {
			t.Errorf("expected ID %d, got %d", created.ID, retrieved.ID)
		}
		if retrieved.Name != "Support/Resistance" {
			t.Errorf("expected name 'Support/Resistance', got %s", retrieved.Name)
		}
		if retrieved.Description != "Trade at S/R levels" {
			t.Errorf("expected description 'Trade at S/R levels', got %s", retrieved.Description)
		}
	})

	t.Run("returns error for non-existent strategy", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("notfound@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Try to get non-existent strategy
		_, err := service.GetStrategy(ctx, 99999, createdUser.ID)

		// Verify error
		if err != ErrStrategyNotFound {
			t.Errorf("expected ErrStrategyNotFound, got %v", err)
		}
	})

	t.Run("user cannot access another user's strategy", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create two users
		user1 := user.NewUser("user1@example.com", "pass1")
		user2 := user.NewUser("user2@example.com", "pass2")

		createdUser1, _ := userRepo.Create(ctx, user1)
		createdUser2, _ := userRepo.Create(ctx, user2)

		// User1 creates a strategy
		req := CreateStrategyRequest{
			Name:        "User1 Strategy",
			Description: "Private to user1",
		}
		strategy, _ := service.CreateStrategy(ctx, createdUser1.ID, req)

		// User2 tries to access user1's strategy
		_, err := service.GetStrategy(ctx, strategy.ID, createdUser2.ID)

		// Should get error (unauthorized)
		if err != ErrStrategyNotFound {
			t.Errorf("expected ErrStrategyNotFound when accessing another user's strategy, got %v", err)
		}
	})
}

func TestStrategyService_GetUserStrategies_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	strategyRepo := persistence.NewStrategyRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(strategyRepo)

	ctx := context.Background()

	t.Run("returns only strategies for the specified user", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create two users
		user1 := user.NewUser("multiuser1@example.com", "pass1")
		user2 := user.NewUser("multiuser2@example.com", "pass2")

		createdUser1, _ := userRepo.Create(ctx, user1)
		createdUser2, _ := userRepo.Create(ctx, user2)

		// Create 3 strategies for user1
		req1 := CreateStrategyRequest{Name: "User1 Strategy 1", Description: "Desc 1"}
		req2 := CreateStrategyRequest{Name: "User1 Strategy 2", Description: "Desc 2"}
		req3 := CreateStrategyRequest{Name: "User1 Strategy 3", Description: "Desc 3"}

		service.CreateStrategy(ctx, createdUser1.ID, req1)
		service.CreateStrategy(ctx, createdUser1.ID, req2)
		service.CreateStrategy(ctx, createdUser1.ID, req3)

		// Create 1 strategy for user2
		req4 := CreateStrategyRequest{Name: "User2 Strategy", Description: "Desc 4"}
		service.CreateStrategy(ctx, createdUser2.ID, req4)

		// Get strategies for user1
		strategies, err := service.GetUserStrategies(ctx, createdUser1.ID)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify user1 has 3 strategies
		if len(strategies) != 3 {
			t.Fatalf("expected 3 strategies for user1, got %d", len(strategies))
		}

		// Verify all strategies belong to user1
		for _, strategy := range strategies {
			if strategy.Name != "User1 Strategy 1" &&
			   strategy.Name != "User1 Strategy 2" &&
			   strategy.Name != "User1 Strategy 3" {
				t.Errorf("unexpected strategy name: %s", strategy.Name)
			}
		}

		// Verify user2 has 1 strategy
		user2Strategies, err := service.GetUserStrategies(ctx, createdUser2.ID)
		if err != nil {
			t.Fatalf("expected no error for user2, got %v", err)
		}
		if len(user2Strategies) != 1 {
			t.Errorf("expected 1 strategy for user2, got %d", len(user2Strategies))
		}
		if user2Strategies[0].Name != "User2 Strategy" {
			t.Errorf("expected 'User2 Strategy', got %s", user2Strategies[0].Name)
		}
	})

	t.Run("returns empty list for user with no strategies", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user with no strategies
		testUser := user.NewUser("nostrategy@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Get strategies
		strategies, err := service.GetUserStrategies(ctx, createdUser.ID)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify empty list
		if len(strategies) != 0 {
			t.Errorf("expected empty list, got %d strategies", len(strategies))
		}
	})
}

func TestStrategyService_UpdateStrategy_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	strategyRepo := persistence.NewStrategyRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(strategyRepo)

	ctx := context.Background()

	t.Run("updates strategy in database", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("updatestrategy@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create strategy
		createReq := CreateStrategyRequest{
			Name:        "Old Name",
			Description: "Old Description",
		}
		created, _ := service.CreateStrategy(ctx, createdUser.ID, createReq)

		// Update strategy
		updateReq := UpdateStrategyRequest{
			Name:        "New Name",
			Description: "New Description",
		}

		updated, err := service.UpdateStrategy(ctx, created.ID, createdUser.ID, updateReq)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify updated fields
		if updated.Name != "New Name" {
			t.Errorf("expected name 'New Name', got %s", updated.Name)
		}
		if updated.Description != "New Description" {
			t.Errorf("expected description 'New Description', got %s", updated.Description)
		}

		// Verify database was updated
		var savedName, savedDesc string
		err = pg.DB.QueryRow("SELECT name, description FROM strategies WHERE id = $1", created.ID).
			Scan(&savedName, &savedDesc)

		if err != nil {
			t.Fatalf("failed to query updated strategy: %v", err)
		}

		if savedName != "New Name" {
			t.Errorf("expected saved name 'New Name', got %s", savedName)
		}
		if savedDesc != "New Description" {
			t.Errorf("expected saved description 'New Description', got %s", savedDesc)
		}
	})

	t.Run("user cannot update another user's strategy", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create two users
		user1 := user.NewUser("updateuser1@example.com", "pass1")
		user2 := user.NewUser("updateuser2@example.com", "pass2")

		createdUser1, _ := userRepo.Create(ctx, user1)
		createdUser2, _ := userRepo.Create(ctx, user2)

		// User1 creates a strategy
		createReq := CreateStrategyRequest{
			Name:        "User1 Strategy",
			Description: "Private",
		}
		strategy, _ := service.CreateStrategy(ctx, createdUser1.ID, createReq)

		// User2 tries to update user1's strategy
		updateReq := UpdateStrategyRequest{
			Name:        "Hacked Name",
			Description: "Hacked",
		}

		_, err := service.UpdateStrategy(ctx, strategy.ID, createdUser2.ID, updateReq)

		// Should get error
		if err != ErrStrategyNotFound {
			t.Errorf("expected ErrStrategyNotFound when updating another user's strategy, got %v", err)
		}

		// Verify original data is unchanged
		var savedName string
		pg.DB.QueryRow("SELECT name FROM strategies WHERE id = $1", strategy.ID).Scan(&savedName)
		if savedName != "User1 Strategy" {
			t.Errorf("strategy should remain unchanged, got name: %s", savedName)
		}
	})
}

func TestStrategyService_DeleteStrategy_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	strategyRepo := persistence.NewStrategyRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(strategyRepo)

	ctx := context.Background()

	t.Run("deletes strategy from database", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("deletestrategy@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create strategy
		req := CreateStrategyRequest{
			Name:        "To Delete",
			Description: "Will be deleted",
		}
		created, _ := service.CreateStrategy(ctx, createdUser.ID, req)

		// Delete strategy
		err := service.DeleteStrategy(ctx, created.ID, createdUser.ID)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify strategy no longer exists in database
		var count int
		err = pg.DB.QueryRow("SELECT COUNT(*) FROM strategies WHERE id = $1", created.ID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query strategy count: %v", err)
		}

		if count != 0 {
			t.Errorf("expected strategy to be deleted, but found %d strategies", count)
		}
	})

	t.Run("user cannot delete another user's strategy", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create two users
		user1 := user.NewUser("deleteuser1@example.com", "pass1")
		user2 := user.NewUser("deleteuser2@example.com", "pass2")

		createdUser1, _ := userRepo.Create(ctx, user1)
		createdUser2, _ := userRepo.Create(ctx, user2)

		// User1 creates a strategy
		req := CreateStrategyRequest{
			Name:        "User1 Strategy",
			Description: "Protected",
		}
		strategy, _ := service.CreateStrategy(ctx, createdUser1.ID, req)

		// User2 tries to delete user1's strategy
		err := service.DeleteStrategy(ctx, strategy.ID, createdUser2.ID)

		// Should get error
		if err != ErrStrategyNotFound {
			t.Errorf("expected ErrStrategyNotFound when deleting another user's strategy, got %v", err)
		}

		// Verify strategy still exists
		var count int
		pg.DB.QueryRow("SELECT COUNT(*) FROM strategies WHERE id = $1", strategy.ID).Scan(&count)
		if count != 1 {
			t.Error("strategy should still exist after failed delete attempt")
		}
	})
}
