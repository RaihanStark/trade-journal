package trade

import (
	"fmt"
	"math"
	"strings"

	tradedom "github.com/raihanstark/trade-journal/internal/domain/trade"
)

// CalculateTradeMetrics calculates pips, P/L, R:R, and status based on trade data
func CalculateTradeMetrics(t *tradedom.Trade) {
	// Calculate R:R for open trades (using take profit) or closed trades (using exit)
	if t.StopLoss != nil {
		var rrFloat float64
		if t.Exit != nil {
			// Trade is closed or has exit - use actual exit price
			rrFloat = calculateRiskReward(t.Type, t.Entry, *t.Exit, *t.StopLoss)
		} else if t.TakeProfit != nil {
			// Trade is open - use take profit for planned R:R
			rrFloat = calculateRiskReward(t.Type, t.Entry, *t.TakeProfit, *t.StopLoss)
		}
		if rrFloat != 0 {
			t.RR = formatRiskReward(rrFloat)
		}
	}

	// Only calculate pips/P/L if we have exit price
	if t.Exit == nil {
		t.Status = tradedom.TradeStatusOpen
		return
	}

	// Calculate pips
	pips := calculatePips(t.Pair, t.Type, t.Entry, *t.Exit)
	t.Pips = &pips

	// Calculate P/L (simplified - assumes $10/pip per lot)
	// In real trading, this would depend on account currency, pair, and lot size
	pl := pips * t.Lots * 10
	t.PL = &pl

	// Set status to closed
	t.Status = tradedom.TradeStatusClosed
}

// calculatePips calculates the pip difference between entry and exit
func calculatePips(pair string, tradeType tradedom.TradeType, entry, exit float64) float64 {
	var pips float64

	// Determine pip multiplier based on pair
	// JPY pairs: 1 pip = 0.01
	// Other pairs: 1 pip = 0.0001
	pipMultiplier := 10000.0
	if strings.Contains(pair, "JPY") {
		pipMultiplier = 100.0
	}

	// Calculate pip difference based on trade type
	if tradeType == tradedom.TradeTypeBuy {
		pips = (exit - entry) * pipMultiplier
	} else if tradeType == tradedom.TradeTypeSell {
		pips = (entry - exit) * pipMultiplier
	}

	return math.Round(pips*100) / 100 // Round to 2 decimal places
}

// calculateRiskReward calculates the risk:reward ratio
func calculateRiskReward(tradeType tradedom.TradeType, entry, exit, stopLoss float64) float64 {
	var risk, reward float64

	if tradeType == tradedom.TradeTypeBuy {
		risk = entry - stopLoss
		reward = exit - entry
	} else if tradeType == tradedom.TradeTypeSell {
		risk = stopLoss - entry
		reward = entry - exit
	}

	if risk == 0 {
		return 0
	}

	rr := reward / risk
	return math.Round(rr*100) / 100 // Round to 2 decimal places
}

// formatRiskReward converts a risk:reward ratio to string format
// Examples: 2.0 -> "1:2", 1.5 -> "1:1.5", 0.5 -> "0.5:1"
func formatRiskReward(ratio float64) string {
	if ratio == 0 {
		return ""
	}

	// Round to 2 decimal places
	ratio = math.Round(ratio*100) / 100

	// If ratio >= 1, format as "1:ratio"
	if ratio >= 1 {
		// Format the ratio part, removing unnecessary trailing zeros
		ratioStr := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", ratio), "0"), ".")
		return fmt.Sprintf("1:%s", ratioStr)
	}

	// If ratio < 1, format as "ratio:1"
	ratioStr := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", ratio), "0"), ".")
	return fmt.Sprintf("%s:1", ratioStr)
}
