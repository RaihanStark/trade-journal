package account

import (
	"context"
	"testing"

	"github.com/raihanstark/trade-journal/internal/domain/user"
	"github.com/raihanstark/trade-journal/internal/infrastructure/persistence"
	"github.com/raihanstark/trade-journal/internal/testutil"
)

func TestAccountService_CreateAccount_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(accountRepo)

	ctx := context.Background()

	t.Run("creates account in database", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		testUser := user.NewUser("accounttest@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		req := CreateAccountRequest{
			Name:          "Demo Account",
			Broker:        "IC Markets",
			AccountNumber: "12345678",
			AccountType:   "demo", // lowercase: database constraint
			Currency:      "USD",
			IsActive:      true,
		}

		result, err := service.CreateAccount(ctx, createdUser.ID, req)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify response
		if result.Name != "Demo Account" {
			t.Errorf("expected name 'Demo Account', got %s", result.Name)
		}
		if result.Broker != "IC Markets" {
			t.Errorf("expected broker 'IC Markets', got %s", result.Broker)
		}
		if result.CurrentBalance != 0.0 {
			t.Errorf("expected initial balance 0, got %.2f", result.CurrentBalance)
		}

		// Verify account was saved to database
		var savedName string
		var savedBalance float64
		err = pg.DB.QueryRow("SELECT name, current_balance FROM accounts WHERE id = $1", result.ID).
			Scan(&savedName, &savedBalance)

		if err != nil {
			t.Fatalf("failed to query saved account: %v", err)
		}

		if savedName != "Demo Account" {
			t.Errorf("expected saved name 'Demo Account', got %s", savedName)
		}
		if savedBalance != 0.0 {
			t.Errorf("expected saved balance 0, got %.2f", savedBalance)
		}
	})
}

func TestAccountRepository_UpdateBalance_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(accountRepo)

	ctx := context.Background()

	t.Run("UpdateBalance adds to existing balance (bug fix verification)", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user using repository
		testUser := user.NewUser("balancetest@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		// Create account using service (cleaner and more maintainable!)
		req := CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, err := service.CreateAccount(ctx, createdUser.ID, req)
		if err != nil {
			t.Fatalf("failed to create test account: %v", err)
		}

		// Verify initial balance is 0
		if account.CurrentBalance != 0.0 {
			t.Fatalf("expected initial balance 0, got %.2f", account.CurrentBalance)
		}

		accountID := account.ID

		// Add $1000 deposit
		_, err = accountRepo.UpdateBalance(ctx, accountID, createdUser.ID, 1000.0)
		if err != nil {
			t.Fatalf("failed to update balance: %v", err)
		}

		// Verify balance is $1000
		var balance1 float64
		err = pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", accountID).Scan(&balance1)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}
		if balance1 != 1000.0 {
			t.Errorf("after deposit of 1000, expected balance 1000, got %.2f", balance1)
		}

		// Add another $500
		_, err = accountRepo.UpdateBalance(ctx, accountID, createdUser.ID, 500.0)
		if err != nil {
			t.Fatalf("failed to update balance second time: %v", err)
		}

		// Verify balance is now $1500 (not $500!)
		var balance2 float64
		err = pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", accountID).Scan(&balance2)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}
		if balance2 != 1500.0 {
			t.Errorf("after adding 500 to 1000, expected balance 1500, got %.2f", balance2)
		}

		// Subtract $200 (withdrawal or loss)
		_, err = accountRepo.UpdateBalance(ctx, accountID, createdUser.ID, -200.0)
		if err != nil {
			t.Fatalf("failed to update balance third time: %v", err)
		}

		// Verify balance is now $1300
		var balance3 float64
		err = pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", accountID).Scan(&balance3)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}
		if balance3 != 1300.0 {
			t.Errorf("after subtracting 200 from 1500, expected balance 1300, got %.2f", balance3)
		}
	})

	t.Run("UpdateBalance handles zero amount", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("zerotest@example.com", "hashedpass")
		createdUser, err := userRepo.Create(ctx, testUser)
		if err != nil {
			t.Fatalf("failed to create test user: %v", err)
		}

		// Create account using service
		req := CreateAccountRequest{
			Name:          "Test Account",
			Broker:        "Test Broker",
			AccountNumber: "123",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		account, err := service.CreateAccount(ctx, createdUser.ID, req)
		if err != nil {
			t.Fatalf("failed to create test account: %v", err)
		}

		accountID := account.ID

		// First, add $100 to have a non-zero balance
		_, err = accountRepo.UpdateBalance(ctx, accountID, createdUser.ID, 100.0)
		if err != nil {
			t.Fatalf("failed to set initial balance: %v", err)
		}

		// Update with zero amount
		_, err = accountRepo.UpdateBalance(ctx, accountID, createdUser.ID, 0.0)
		if err != nil {
			t.Fatalf("failed to update balance with zero: %v", err)
		}

		// Balance should still be 100
		var balance float64
		err = pg.DB.QueryRow("SELECT current_balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
		if err != nil {
			t.Fatalf("failed to query balance: %v", err)
		}
		if balance != 100.0 {
			t.Errorf("after adding 0, expected balance 100, got %.2f", balance)
		}
	})
}

func TestAccountService_GetUserAccounts_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(accountRepo)

	ctx := context.Background()

	t.Run("returns only accounts for the specified user", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create two users
		user1 := user.NewUser("user1@example.com", "pass1")
		user2 := user.NewUser("user2@example.com", "pass2")

		createdUser1, _ := userRepo.Create(ctx, user1)
		createdUser2, _ := userRepo.Create(ctx, user2)

		// Create 2 accounts for user1
		req1 := CreateAccountRequest{
			Name:          "User1 Account 1",
			Broker:        "Broker A",
			AccountNumber: "111",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		req2 := CreateAccountRequest{
			Name:          "User1 Account 2",
			Broker:        "Broker B",
			AccountNumber: "222",
			AccountType:   "live",
			Currency:      "EUR",
			IsActive:      true,
		}

		service.CreateAccount(ctx, createdUser1.ID, req1)
		service.CreateAccount(ctx, createdUser1.ID, req2)

		// Create 1 account for user2
		req3 := CreateAccountRequest{
			Name:          "User2 Account",
			Broker:        "Broker C",
			AccountNumber: "333",
			AccountType:   "demo",
			Currency:      "GBP",
			IsActive:      true,
		}
		service.CreateAccount(ctx, createdUser2.ID, req3)

		// Get accounts for user1
		accounts, err := service.GetUserAccounts(ctx, createdUser1.ID)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify user1 has 2 accounts
		if len(accounts) != 2 {
			t.Fatalf("expected 2 accounts for user1, got %d", len(accounts))
		}

		// Verify both accounts belong to user1
		for _, acc := range accounts {
			if acc.Name != "User1 Account 1" && acc.Name != "User1 Account 2" {
				t.Errorf("unexpected account name: %s", acc.Name)
			}
		}

		// Verify user2 has 1 account
		user2Accounts, err := service.GetUserAccounts(ctx, createdUser2.ID)
		if err != nil {
			t.Fatalf("expected no error for user2, got %v", err)
		}
		if len(user2Accounts) != 1 {
			t.Errorf("expected 1 account for user2, got %d", len(user2Accounts))
		}
		if user2Accounts[0].Name != "User2 Account" {
			t.Errorf("expected 'User2 Account', got %s", user2Accounts[0].Name)
		}
	})
}

func TestAccountService_UpdateAccount_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(accountRepo)

	ctx := context.Background()

	t.Run("updates account in database", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("updatetest@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account
		createReq := CreateAccountRequest{
			Name:          "Old Name",
			Broker:        "Old Broker",
			AccountNumber: "111",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		created, _ := service.CreateAccount(ctx, createdUser.ID, createReq)

		// Update account
		updateReq := UpdateAccountRequest{
			Name:          "New Name",
			Broker:        "New Broker",
			AccountNumber: "222",
			AccountType:   "live",
			Currency:      "EUR",
			IsActive:      false,
		}

		updated, err := service.UpdateAccount(ctx, created.ID, createdUser.ID, updateReq)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify updated fields
		if updated.Name != "New Name" {
			t.Errorf("expected name 'New Name', got %s", updated.Name)
		}
		if updated.Broker != "New Broker" {
			t.Errorf("expected broker 'New Broker', got %s", updated.Broker)
		}
		if updated.AccountType != "live" {
			t.Errorf("expected type 'live', got %s", updated.AccountType)
		}

		// Verify database was updated
		var savedName, savedBroker string
		err = pg.DB.QueryRow("SELECT name, broker FROM accounts WHERE id = $1", created.ID).
			Scan(&savedName, &savedBroker)

		if err != nil {
			t.Fatalf("failed to query updated account: %v", err)
		}

		if savedName != "New Name" {
			t.Errorf("expected saved name 'New Name', got %s", savedName)
		}
		if savedBroker != "New Broker" {
			t.Errorf("expected saved broker 'New Broker', got %s", savedBroker)
		}
	})
}

func TestAccountService_DeleteAccount_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	pg := testutil.SetupTestDatabase(t)
	accountRepo := persistence.NewAccountRepository(pg.Queries)
	userRepo := persistence.NewUserRepository(pg.Queries)
	service := NewService(accountRepo)

	ctx := context.Background()

	t.Run("deletes account from database", func(t *testing.T) {
		testutil.TruncateTables(t, pg.DB)

		// Create user
		testUser := user.NewUser("deletetest@example.com", "hashedpass")
		createdUser, _ := userRepo.Create(ctx, testUser)

		// Create account
		req := CreateAccountRequest{
			Name:          "To Delete",
			Broker:        "Test Broker",
			AccountNumber: "999",
			AccountType:   "demo",
			Currency:      "USD",
			IsActive:      true,
		}
		created, _ := service.CreateAccount(ctx, createdUser.ID, req)

		// Delete account
		err := service.DeleteAccount(ctx, created.ID, createdUser.ID)

		// Verify no error
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify account no longer exists in database
		var count int
		err = pg.DB.QueryRow("SELECT COUNT(*) FROM accounts WHERE id = $1", created.ID).Scan(&count)
		if err != nil {
			t.Fatalf("failed to query account count: %v", err)
		}

		if count != 0 {
			t.Errorf("expected account to be deleted, but found %d accounts", count)
		}
	})
}
