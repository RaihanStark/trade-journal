package analytics

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/raihanstark/trade-journal/internal/db"
	"github.com/raihanstark/trade-journal/internal/domain/analytics"
)

// Calculator handles all analytics calculations
type Calculator struct{}

// NewCalculator creates a new Calculator instance
func NewCalculator() *Calculator {
	return &Calculator{}
}

// CalculateAnalytics performs all analytics calculations on the given trades
func (c *Calculator) CalculateAnalytics(trades []db.Trade) *analytics.Analytics {
	result := &analytics.Analytics{}

	// Filter only BUY and SELL trades with P/L (closed trades)
	closedTrades := c.filterClosedTrades(trades)

	if len(closedTrades) == 0 {
		return result
	}

	// Calculate basic metrics
	totalPL, winningTrades, losingTrades, totalWinPL, totalLossPL, largestWin, largestLoss := c.calculateBasicMetrics(closedTrades)

	result.TotalTrades = int64(len(closedTrades))
	result.WinningTrades = winningTrades
	result.LosingTrades = losingTrades
	result.TotalPL = totalPL
	result.LargestWin = largestWin
	result.LargestLoss = largestLoss

	// Calculate win rate
	if result.TotalTrades > 0 {
		result.WinRate = (float64(winningTrades) / float64(result.TotalTrades)) * 100
	}

	// Calculate average win and loss
	if winningTrades > 0 {
		result.AvgWin = totalWinPL / float64(winningTrades)
	}
	if losingTrades > 0 {
		result.AvgLoss = -totalLossPL / float64(losingTrades)
	}

	// Calculate profit factor
	if totalLossPL > 0 {
		result.ProfitFactor = totalWinPL / totalLossPL
	}

	// Calculate streaks
	result.ConsecutiveWins, result.ConsecutiveLosses, result.BestStreak, result.WorstStreak = c.calculateStreaks(closedTrades)

	// Calculate Sharpe Ratio (simplified version)
	result.SharpeRatio = c.calculateSharpeRatio(closedTrades)

	// Calculate Max Drawdown
	result.MaxDrawdown = c.calculateMaxDrawdown(closedTrades)

	return result
}

// filterClosedTrades filters only BUY and SELL trades with valid P/L
func (c *Calculator) filterClosedTrades(trades []db.Trade) []db.Trade {
	var closedTrades []db.Trade
	for _, trade := range trades {
		if (trade.Type == db.TradeTypeBUY || trade.Type == db.TradeTypeSELL) && trade.Pl.Valid {
			closedTrades = append(closedTrades, trade)
		}
	}
	return closedTrades
}

// calculateBasicMetrics calculates basic P/L metrics
func (c *Calculator) calculateBasicMetrics(trades []db.Trade) (
	totalPL float64,
	winningTrades, losingTrades int64,
	totalWinPL, totalLossPL float64,
	largestWin, largestLoss float64,
) {
	for _, trade := range trades {
		pl := parseFloatFromNullString(trade.Pl)
		totalPL += pl

		if pl > 0 {
			winningTrades++
			totalWinPL += pl
			if pl > largestWin {
				largestWin = pl
			}
		} else if pl < 0 {
			losingTrades++
			totalLossPL += math.Abs(pl)
			if pl < largestLoss {
				largestLoss = pl
			}
		}
	}
	return
}

// calculateStreaks calculates current and best/worst streaks
func (c *Calculator) calculateStreaks(trades []db.Trade) (currentWins, currentLosses, bestStreak, worstStreak int64) {
	if len(trades) == 0 {
		return 0, 0, 0, 0
	}

	var currentStreak int64
	var maxWinStreak, maxLossStreak int64

	for _, trade := range trades {
		pl := parseFloatFromNullString(trade.Pl)

		if pl > 0 {
			if currentStreak > 0 {
				currentStreak++
			} else {
				currentStreak = 1
			}
			if currentStreak > maxWinStreak {
				maxWinStreak = currentStreak
			}
		} else if pl < 0 {
			if currentStreak < 0 {
				currentStreak--
			} else {
				currentStreak = -1
			}
			if currentStreak < maxLossStreak {
				maxLossStreak = currentStreak
			}
		}
	}

	// Current streak
	if currentStreak > 0 {
		currentWins = currentStreak
	} else if currentStreak < 0 {
		currentLosses = -currentStreak
	}

	bestStreak = maxWinStreak
	worstStreak = -maxLossStreak

	return
}

// calculateSharpeRatio calculates the Sharpe Ratio (simplified, assuming risk-free rate = 0)
func (c *Calculator) calculateSharpeRatio(trades []db.Trade) float64 {
	if len(trades) < 2 {
		return 0
	}

	// Calculate average return and standard deviation
	var returns []float64
	var sum float64

	for _, trade := range trades {
		pl := parseFloatFromNullString(trade.Pl)
		returns = append(returns, pl)
		sum += pl
	}

	avgReturn := sum / float64(len(returns))

	// Calculate standard deviation
	var variance float64
	for _, ret := range returns {
		variance += math.Pow(ret-avgReturn, 2)
	}
	variance /= float64(len(returns))
	stdDev := math.Sqrt(variance)

	if stdDev == 0 {
		return 0
	}

	// Simplified Sharpe Ratio (assuming risk-free rate = 0)
	return avgReturn / stdDev
}

// calculateMaxDrawdown calculates the maximum drawdown
func (c *Calculator) calculateMaxDrawdown(trades []db.Trade) float64 {
	if len(trades) == 0 {
		return 0
	}

	var peak float64
	var maxDrawdown float64
	var runningBalance float64

	for _, trade := range trades {
		pl := parseFloatFromNullString(trade.Pl)
		runningBalance += pl

		if runningBalance > peak {
			peak = runningBalance
		}

		drawdown := peak - runningBalance
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return -maxDrawdown
}

// parseFloatFromNullString converts sql.NullString to float64
func parseFloatFromNullString(ns sql.NullString) float64 {
	if !ns.Valid || ns.String == "" {
		return 0
	}

	var f float64
	_, err := fmt.Sscanf(ns.String, "%f", &f)
	if err != nil {
		return 0
	}
	return f
}
