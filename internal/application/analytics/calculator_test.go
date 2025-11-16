package analytics

import (
	"database/sql"
	"math"
	"testing"

	"github.com/raihanstark/trade-journal/internal/db"
	domain "github.com/raihanstark/trade-journal/internal/domain/analytics"
)

// Helper function to create a NullString
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func TestParseFloatFromNullString(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected float64
	}{
		{
			name:     "valid positive number",
			input:    nullString("123.45"),
			expected: 123.45,
		},
		{
			name:     "valid negative number",
			input:    nullString("-67.89"),
			expected: -67.89,
		},
		{
			name:     "zero",
			input:    nullString("0"),
			expected: 0,
		},
		{
			name:     "invalid null string",
			input:    sql.NullString{Valid: false},
			expected: 0,
		},
		{
			name:     "empty string",
			input:    nullString(""),
			expected: 0,
		},
		{
			name:     "invalid format",
			input:    nullString("abc"),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseFloatFromNullString(tt.input)
			if result != tt.expected {
				t.Errorf("parseFloatFromNullString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFilterClosedTrades(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		trades   []db.Trade
		expected int
	}{
		{
			name: "only closed trades",
			trades: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("100.50")},
				{Type: db.TradeTypeSELL, Pl: nullString("-50.25")},
			},
			expected: 2,
		},
		{
			name: "mixed with open trades",
			trades: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("100.50")},
				{Type: db.TradeTypeSELL, Pl: sql.NullString{Valid: false}}, // open trade
				{Type: db.TradeTypeSELL, Pl: nullString("-50.25")},
			},
			expected: 2,
		},
		{
			name: "mixed with deposit/withdraw",
			trades: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("100.50")},
				{Type: db.TradeTypeDEPOSIT, Pl: nullString("1000")}, // should be filtered out
				{Type: db.TradeTypeWITHDRAW, Pl: nullString("-500")}, // should be filtered out
				{Type: db.TradeTypeSELL, Pl: nullString("-50.25")},
			},
			expected: 2,
		},
		{
			name:     "empty trades",
			trades:   []db.Trade{},
			expected: 0,
		},
		{
			name: "only open trades",
			trades: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: sql.NullString{Valid: false}},
				{Type: db.TradeTypeSELL, Pl: sql.NullString{Valid: false}},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.filterClosedTrades(tt.trades)
			if len(result) != tt.expected {
				t.Errorf("filterClosedTrades() returned %d trades, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestCalculateBasicMetrics(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name              string
		trades            []db.Trade
		expectedTotalPL   float64
		expectedWinning   int64
		expectedLosing    int64
		expectedTotalWin  float64
		expectedTotalLoss float64
		expectedLargestW  float64
		expectedLargestL  float64
	}{
		{
			name: "mixed wins and losses",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("-50")},
				{Pl: nullString("200")},
				{Pl: nullString("-75")},
			},
			expectedTotalPL:   175,
			expectedWinning:   2,
			expectedLosing:    2,
			expectedTotalWin:  300,
			expectedTotalLoss: 125,
			expectedLargestW:  200,
			expectedLargestL:  -75,
		},
		{
			name: "only wins",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("200")},
				{Pl: nullString("50")},
			},
			expectedTotalPL:   350,
			expectedWinning:   3,
			expectedLosing:    0,
			expectedTotalWin:  350,
			expectedTotalLoss: 0,
			expectedLargestW:  200,
			expectedLargestL:  0,
		},
		{
			name: "only losses",
			trades: []db.Trade{
				{Pl: nullString("-100")},
				{Pl: nullString("-200")},
				{Pl: nullString("-50")},
			},
			expectedTotalPL:   -350,
			expectedWinning:   0,
			expectedLosing:    3,
			expectedTotalWin:  0,
			expectedTotalLoss: 350,
			expectedLargestW:  0,
			expectedLargestL:  -200,
		},
		{
			name: "with zero P/L trades",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("0")},
				{Pl: nullString("-50")},
			},
			expectedTotalPL:   50,
			expectedWinning:   1,
			expectedLosing:    1,
			expectedTotalWin:  100,
			expectedTotalLoss: 50,
			expectedLargestW:  100,
			expectedLargestL:  -50,
		},
		{
			name:              "empty trades",
			trades:            []db.Trade{},
			expectedTotalPL:   0,
			expectedWinning:   0,
			expectedLosing:    0,
			expectedTotalWin:  0,
			expectedTotalLoss: 0,
			expectedLargestW:  0,
			expectedLargestL:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totalPL, winning, losing, totalWin, totalLoss, largestW, largestL := calc.calculateBasicMetrics(tt.trades)

			if totalPL != tt.expectedTotalPL {
				t.Errorf("totalPL = %v, want %v", totalPL, tt.expectedTotalPL)
			}
			if winning != tt.expectedWinning {
				t.Errorf("winning = %v, want %v", winning, tt.expectedWinning)
			}
			if losing != tt.expectedLosing {
				t.Errorf("losing = %v, want %v", losing, tt.expectedLosing)
			}
			if totalWin != tt.expectedTotalWin {
				t.Errorf("totalWin = %v, want %v", totalWin, tt.expectedTotalWin)
			}
			if totalLoss != tt.expectedTotalLoss {
				t.Errorf("totalLoss = %v, want %v", totalLoss, tt.expectedTotalLoss)
			}
			if largestW != tt.expectedLargestW {
				t.Errorf("largestWin = %v, want %v", largestW, tt.expectedLargestW)
			}
			if largestL != tt.expectedLargestL {
				t.Errorf("largestLoss = %v, want %v", largestL, tt.expectedLargestL)
			}
		})
	}
}

func TestCalculateStreaks(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name           string
		trades         []db.Trade
		expectedCurWin int64
		expectedCurLos int64
		expectedBest   int64
		expectedWorst  int64
	}{
		{
			name: "alternating wins and losses",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("-50")},
				{Pl: nullString("75")},
				{Pl: nullString("-25")},
			},
			expectedCurWin: 0,
			expectedCurLos: 1,
			expectedBest:   1,
			expectedWorst:  1,
		},
		{
			name: "consecutive wins",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("50")},
				{Pl: nullString("75")},
				{Pl: nullString("25")},
			},
			expectedCurWin: 4,
			expectedCurLos: 0,
			expectedBest:   4,
			expectedWorst:  0,
		},
		{
			name: "consecutive losses",
			trades: []db.Trade{
				{Pl: nullString("-100")},
				{Pl: nullString("-50")},
				{Pl: nullString("-75")},
			},
			expectedCurWin: 0,
			expectedCurLos: 3,
			expectedBest:   0,
			expectedWorst:  3,
		},
		{
			name: "win streak then loss streak",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("50")},
				{Pl: nullString("75")},
				{Pl: nullString("-25")},
				{Pl: nullString("-50")},
			},
			expectedCurWin: 0,
			expectedCurLos: 2,
			expectedBest:   3,
			expectedWorst:  2,
		},
		{
			name: "loss streak then win streak",
			trades: []db.Trade{
				{Pl: nullString("-100")},
				{Pl: nullString("-50")},
				{Pl: nullString("-75")},
				{Pl: nullString("25")},
				{Pl: nullString("50")},
				{Pl: nullString("100")},
			},
			expectedCurWin: 3,
			expectedCurLos: 0,
			expectedBest:   3,
			expectedWorst:  3,
		},
		{
			name: "with zero P/L trades",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("0")}, // Zero P/L doesn't break streak
				{Pl: nullString("50")},
			},
			expectedCurWin: 2, // Zero doesn't break the win streak
			expectedCurLos: 0,
			expectedBest:   2,
			expectedWorst:  0,
		},
		{
			name:           "empty trades",
			trades:         []db.Trade{},
			expectedCurWin: 0,
			expectedCurLos: 0,
			expectedBest:   0,
			expectedWorst:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curWin, curLos, best, worst := calc.calculateStreaks(tt.trades)

			if curWin != tt.expectedCurWin {
				t.Errorf("currentWins = %v, want %v", curWin, tt.expectedCurWin)
			}
			if curLos != tt.expectedCurLos {
				t.Errorf("currentLosses = %v, want %v", curLos, tt.expectedCurLos)
			}
			if best != tt.expectedBest {
				t.Errorf("bestStreak = %v, want %v", best, tt.expectedBest)
			}
			if worst != tt.expectedWorst {
				t.Errorf("worstStreak = %v, want %v", worst, tt.expectedWorst)
			}
		})
	}
}

func TestCalculateSharpeRatio(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		trades   []db.Trade
		expected float64
	}{
		{
			name: "consistent positive returns",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("100")},
				{Pl: nullString("100")},
				{Pl: nullString("100")},
			},
			expected: 0, // stdDev = 0, so Sharpe = 0
		},
		{
			name: "mixed returns",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("-50")},
				{Pl: nullString("200")},
				{Pl: nullString("-75")},
			},
			// Average = (100 - 50 + 200 - 75) / 4 = 43.75
			// Variance = ((56.25)² + (-93.75)² + (156.25)² + (-118.75)²) / 4
			// Variance = (3164.0625 + 8789.0625 + 24414.0625 + 14101.5625) / 4 = 12617.1875
			// StdDev = sqrt(12617.1875) = 112.33
			// Sharpe = 43.75 / 112.33 ≈ 0.3895
			expected: 0.3895, // approximately
		},
		{
			name: "only one trade",
			trades: []db.Trade{
				{Pl: nullString("100")},
			},
			expected: 0, // Less than 2 trades
		},
		{
			name:     "empty trades",
			trades:   []db.Trade{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.calculateSharpeRatio(tt.trades)

			// For non-zero expectations, use approximate comparison
			if tt.expected != 0 {
				if math.Abs(result-tt.expected) > 0.01 {
					t.Errorf("calculateSharpeRatio() = %v, want approximately %v", result, tt.expected)
				}
			} else {
				if result != tt.expected {
					t.Errorf("calculateSharpeRatio() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestCalculateMaxDrawdown(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		trades   []db.Trade
		expected float64
	}{
		{
			name: "continuous profit",
			trades: []db.Trade{
				{Pl: nullString("100")},
				{Pl: nullString("50")},
				{Pl: nullString("75")},
			},
			expected: 0, // No drawdown
		},
		{
			name: "drawdown after peak",
			trades: []db.Trade{
				{Pl: nullString("100")}, // Balance: 100, Peak: 100
				{Pl: nullString("50")},  // Balance: 150, Peak: 150
				{Pl: nullString("-75")}, // Balance: 75, Peak: 150, Drawdown: 75
				{Pl: nullString("-25")}, // Balance: 50, Peak: 150, Drawdown: 100
			},
			expected: -100,
		},
		{
			name: "recovery after drawdown",
			trades: []db.Trade{
				{Pl: nullString("200")},  // Balance: 200, Peak: 200
				{Pl: nullString("-100")}, // Balance: 100, Peak: 200, Drawdown: 100
				{Pl: nullString("-50")},  // Balance: 50, Peak: 200, Drawdown: 150 (max)
				{Pl: nullString("150")},  // Balance: 200, Peak: 200, Drawdown: 0
				{Pl: nullString("50")},   // Balance: 250, Peak: 250
			},
			expected: -150,
		},
		{
			name: "only losses from zero",
			trades: []db.Trade{
				{Pl: nullString("-50")},  // Balance: -50, Peak: 0, DD: 50
				{Pl: nullString("-25")},  // Balance: -75, Peak: 0, DD: 75
				{Pl: nullString("-75")},  // Balance: -150, Peak: 0, DD: 150
			},
			expected: -150, // Drawdown from peak of 0 to -150
		},
		{
			name:     "empty trades",
			trades:   []db.Trade{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.calculateMaxDrawdown(tt.trades)

			if result != tt.expected {
				t.Errorf("calculateMaxDrawdown() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateAnalytics(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		trades   []db.Trade
		validate func(t *testing.T, result *domain.Analytics)
	}{
		{
			name: "comprehensive analytics test",
			trades: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("100")},
				{Type: db.TradeTypeSELL, Pl: nullString("-50")},
				{Type: db.TradeTypeBUY, Pl: nullString("200")},
				{Type: db.TradeTypeSELL, Pl: nullString("-75")},
				{Type: db.TradeTypeBUY, Pl: nullString("150")},
			},
			validate: func(t *testing.T, result *domain.Analytics) {
				// Total P/L: 100 - 50 + 200 - 75 + 150 = 325
				if result.TotalPL != 325 {
					t.Errorf("TotalPL = %v, want 325", result.TotalPL)
				}

				// Total trades: 5
				if result.TotalTrades != 5 {
					t.Errorf("TotalTrades = %v, want 5", result.TotalTrades)
				}

				// Winning trades: 3, Losing trades: 2
				if result.WinningTrades != 3 {
					t.Errorf("WinningTrades = %v, want 3", result.WinningTrades)
				}
				if result.LosingTrades != 2 {
					t.Errorf("LosingTrades = %v, want 2", result.LosingTrades)
				}

				// Win rate: 3/5 * 100 = 60%
				if result.WinRate != 60 {
					t.Errorf("WinRate = %v, want 60", result.WinRate)
				}

				// Avg win: (100 + 200 + 150) / 3 = 150
				if result.AvgWin != 150 {
					t.Errorf("AvgWin = %v, want 150", result.AvgWin)
				}

				// Avg loss: -(50 + 75) / 2 = -62.5
				if result.AvgLoss != -62.5 {
					t.Errorf("AvgLoss = %v, want -62.5", result.AvgLoss)
				}

				// Profit factor: 450 / 125 = 3.6
				if result.ProfitFactor != 3.6 {
					t.Errorf("ProfitFactor = %v, want 3.6", result.ProfitFactor)
				}

				// Largest win: 200
				if result.LargestWin != 200 {
					t.Errorf("LargestWin = %v, want 200", result.LargestWin)
				}

				// Largest loss: -75
				if result.LargestLoss != -75 {
					t.Errorf("LargestLoss = %v, want -75", result.LargestLoss)
				}
			},
		},
		{
			name: "filters out open trades and deposits/withdrawals",
			trades: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: nullString("100")},
				{Type: db.TradeTypeSELL, Pl: sql.NullString{Valid: false}}, // open trade
				{Type: db.TradeTypeDEPOSIT, Pl: nullString("1000")},        // deposit
				{Type: db.TradeTypeWITHDRAW, Pl: nullString("-500")},       // withdraw
				{Type: db.TradeTypeBUY, Pl: nullString("50")},
			},
			validate: func(t *testing.T, result *domain.Analytics) {
				// Should only count 2 closed BUY/SELL trades
				if result.TotalTrades != 2 {
					t.Errorf("TotalTrades = %v, want 2 (should filter out open trades and deposits/withdrawals)", result.TotalTrades)
				}

				if result.TotalPL != 150 {
					t.Errorf("TotalPL = %v, want 150", result.TotalPL)
				}
			},
		},
		{
			name:   "empty trades",
			trades: []db.Trade{},
			validate: func(t *testing.T, result *domain.Analytics) {
				if result.TotalTrades != 0 {
					t.Errorf("TotalTrades = %v, want 0", result.TotalTrades)
				}
				if result.TotalPL != 0 {
					t.Errorf("TotalPL = %v, want 0", result.TotalPL)
				}
				if result.WinRate != 0 {
					t.Errorf("WinRate = %v, want 0", result.WinRate)
				}
			},
		},
		{
			name: "only open trades",
			trades: []db.Trade{
				{Type: db.TradeTypeBUY, Pl: sql.NullString{Valid: false}},
				{Type: db.TradeTypeSELL, Pl: sql.NullString{Valid: false}},
			},
			validate: func(t *testing.T, result *domain.Analytics) {
				if result.TotalTrades != 0 {
					t.Errorf("TotalTrades = %v, want 0", result.TotalTrades)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.CalculateAnalytics(tt.trades)
			tt.validate(t, result)
		})
	}
}
