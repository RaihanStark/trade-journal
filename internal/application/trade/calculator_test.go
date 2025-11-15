package trade

import (
	"testing"

	tradedom "github.com/raihanstark/trade-journal/internal/domain/trade"
)

func TestCalculatePips(t *testing.T) {
	tests := []struct {
		name      string
		pair      string
		tradeType tradedom.TradeType
		entry     float64
		exit      float64
		want      float64
	}{
		{
			name:      "BUY trade - EUR/USD - profit",
			pair:      "EUR/USD",
			tradeType: tradedom.TradeTypeBuy,
			entry:     1.1000,
			exit:      1.1050,
			want:      50.0,
		},
		{
			name:      "BUY trade - EUR/USD - loss",
			pair:      "EUR/USD",
			tradeType: tradedom.TradeTypeBuy,
			entry:     1.1000,
			exit:      1.0950,
			want:      -50.0,
		},
		{
			name:      "SELL trade - EUR/USD - profit",
			pair:      "EUR/USD",
			tradeType: tradedom.TradeTypeSell,
			entry:     1.1000,
			exit:      1.0950,
			want:      50.0,
		},
		{
			name:      "SELL trade - EUR/USD - loss",
			pair:      "EUR/USD",
			tradeType: tradedom.TradeTypeSell,
			entry:     1.1000,
			exit:      1.1050,
			want:      -50.0,
		},
		{
			name:      "BUY trade - USD/JPY - profit",
			pair:      "USD/JPY",
			tradeType: tradedom.TradeTypeBuy,
			entry:     110.00,
			exit:      110.50,
			want:      50.0,
		},
		{
			name:      "BUY trade - USD/JPY - loss",
			pair:      "USD/JPY",
			tradeType: tradedom.TradeTypeBuy,
			entry:     110.00,
			exit:      109.50,
			want:      -50.0,
		},
		{
			name:      "SELL trade - USD/JPY - profit",
			pair:      "USD/JPY",
			tradeType: tradedom.TradeTypeSell,
			entry:     110.00,
			exit:      109.50,
			want:      50.0,
		},
		{
			name:      "BUY trade - GBP/JPY - profit",
			pair:      "GBP/JPY",
			tradeType: tradedom.TradeTypeBuy,
			entry:     150.00,
			exit:      150.75,
			want:      75.0,
		},
		{
			name:      "BUY trade - small pip movement",
			pair:      "EUR/USD",
			tradeType: tradedom.TradeTypeBuy,
			entry:     1.1000,
			exit:      1.1001,
			want:      1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculatePips(tt.pair, tt.tradeType, tt.entry, tt.exit)
			if got != tt.want {
				t.Errorf("calculatePips() = %.2f, want %.2f", got, tt.want)
			}
		})
	}
}

func TestCalculateRiskReward(t *testing.T) {
	tests := []struct {
		name      string
		tradeType tradedom.TradeType
		entry     float64
		exit      float64
		stopLoss  float64
		want      float64
	}{
		{
			name:      "BUY trade - 1:2 risk reward",
			tradeType: tradedom.TradeTypeBuy,
			entry:     1.1000,
			exit:      1.1040,
			stopLoss:  1.0980,
			want:      2.0,
		},
		{
			name:      "BUY trade - 1:1 risk reward",
			tradeType: tradedom.TradeTypeBuy,
			entry:     1.1000,
			exit:      1.1020,
			stopLoss:  1.0980,
			want:      1.0,
		},
		{
			name:      "BUY trade - 1:3 risk reward",
			tradeType: tradedom.TradeTypeBuy,
			entry:     1.1000,
			exit:      1.1060,
			stopLoss:  1.0980,
			want:      3.0,
		},
		{
			name:      "BUY trade - stopped out (negative)",
			tradeType: tradedom.TradeTypeBuy,
			entry:     1.1000,
			exit:      1.0990,
			stopLoss:  1.0980,
			want:      -0.5,
		},
		{
			name:      "SELL trade - 1:2 risk reward",
			tradeType: tradedom.TradeTypeSell,
			entry:     1.1000,
			exit:      1.0960,
			stopLoss:  1.1020,
			want:      2.0,
		},
		{
			name:      "SELL trade - stopped out (negative)",
			tradeType: tradedom.TradeTypeSell,
			entry:     1.1000,
			exit:      1.1010,
			stopLoss:  1.1020,
			want:      -0.5,
		},
		{
			name:      "BUY trade - zero risk (should return 0)",
			tradeType: tradedom.TradeTypeBuy,
			entry:     1.1000,
			exit:      1.1020,
			stopLoss:  1.1000,
			want:      0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateRiskReward(tt.tradeType, tt.entry, tt.exit, tt.stopLoss)
			if got != tt.want {
				t.Errorf("calculateRiskReward() = %.2f, want %.2f", got, tt.want)
			}
		})
	}
}

func TestFormatRiskReward(t *testing.T) {
	tests := []struct {
		name  string
		ratio float64
		want  string
	}{
		{
			name:  "ratio 2.0",
			ratio: 2.0,
			want:  "1:2",
		},
		{
			name:  "ratio 1.5",
			ratio: 1.5,
			want:  "1:1.5",
		},
		{
			name:  "ratio 1.0",
			ratio: 1.0,
			want:  "1:1",
		},
		{
			name:  "ratio 0.5",
			ratio: 0.5,
			want:  "0.5:1",
		},
		{
			name:  "ratio 0.75",
			ratio: 0.75,
			want:  "0.75:1",
		},
		{
			name:  "ratio 3.25",
			ratio: 3.25,
			want:  "1:3.25",
		},
		{
			name:  "ratio 0.0",
			ratio: 0.0,
			want:  "",
		},
		{
			name:  "negative ratio -0.5",
			ratio: -0.5,
			want:  "-0.5:1",
		},
		{
			name:  "ratio with trailing zeros 2.00",
			ratio: 2.00,
			want:  "1:2",
		},
		{
			name:  "ratio 1.10",
			ratio: 1.10,
			want:  "1:1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatRiskReward(tt.ratio)
			if got != tt.want {
				t.Errorf("formatRiskReward() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCalculateTradeMetrics(t *testing.T) {
	t.Run("BUY trade - closed with profit", func(t *testing.T) {
		exit := 1.1050
		trade := &tradedom.Trade{
			Pair:       "EUR/USD",
			Type:       tradedom.TradeTypeBuy,
			Entry:      1.1000,
			Exit:       &exit,
			Lots:       1.0,
			StopLoss:   floatPtr(1.0980),
			TakeProfit: floatPtr(1.1060),
		}

		CalculateTradeMetrics(trade)

		// Assert pips
		if trade.Pips == nil {
			t.Fatal("expected pips to be calculated, got nil")
		}
		if *trade.Pips != 50.0 {
			t.Errorf("expected pips 50.0, got %.2f", *trade.Pips)
		}

		// Assert P/L
		if trade.PL == nil {
			t.Fatal("expected PL to be calculated, got nil")
		}
		expectedPL := 500.0 // 50 pips * 1 lot * $10
		if *trade.PL != expectedPL {
			t.Errorf("expected PL %.2f, got %.2f", expectedPL, *trade.PL)
		}

		// Assert R:R (using actual exit)
		if trade.RR != "1:2.5" {
			t.Errorf("expected RR '1:2.5', got %s", trade.RR)
		}

		// Assert status
		if trade.Status != tradedom.TradeStatusClosed {
			t.Errorf("expected status CLOSED, got %s", trade.Status)
		}
	})

	t.Run("BUY trade - open with take profit", func(t *testing.T) {
		trade := &tradedom.Trade{
			Pair:       "EUR/USD",
			Type:       tradedom.TradeTypeBuy,
			Entry:      1.1000,
			Exit:       nil, // Open trade
			Lots:       1.0,
			StopLoss:   floatPtr(1.0980),
			TakeProfit: floatPtr(1.1060),
		}

		CalculateTradeMetrics(trade)

		// Assert pips and P/L are nil for open trades
		if trade.Pips != nil {
			t.Errorf("expected pips to be nil for open trade, got %.2f", *trade.Pips)
		}
		if trade.PL != nil {
			t.Errorf("expected PL to be nil for open trade, got %.2f", *trade.PL)
		}

		// Assert R:R (using take profit)
		if trade.RR != "1:3" {
			t.Errorf("expected RR '1:3', got %s", trade.RR)
		}

		// Assert status
		if trade.Status != tradedom.TradeStatusOpen {
			t.Errorf("expected status OPEN, got %s", trade.Status)
		}
	})

	t.Run("SELL trade - closed with loss", func(t *testing.T) {
		exit := 1.1050
		trade := &tradedom.Trade{
			Pair:       "EUR/USD",
			Type:       tradedom.TradeTypeSell,
			Entry:      1.1000,
			Exit:       &exit,
			Lots:       0.5,
			StopLoss:   floatPtr(1.1020),
			TakeProfit: floatPtr(1.0950),
		}

		CalculateTradeMetrics(trade)

		// Assert pips (negative for loss)
		if trade.Pips == nil {
			t.Fatal("expected pips to be calculated, got nil")
		}
		if *trade.Pips != -50.0 {
			t.Errorf("expected pips -50.0, got %.2f", *trade.Pips)
		}

		// Assert P/L (negative)
		if trade.PL == nil {
			t.Fatal("expected PL to be calculated, got nil")
		}
		expectedPL := -250.0 // -50 pips * 0.5 lot * $10
		if *trade.PL != expectedPL {
			t.Errorf("expected PL %.2f, got %.2f", expectedPL, *trade.PL)
		}

		// Assert status
		if trade.Status != tradedom.TradeStatusClosed {
			t.Errorf("expected status CLOSED, got %s", trade.Status)
		}
	})

	t.Run("trade without stop loss - no R:R calculation", func(t *testing.T) {
		exit := 1.1050
		trade := &tradedom.Trade{
			Pair:  "EUR/USD",
			Type:  tradedom.TradeTypeBuy,
			Entry: 1.1000,
			Exit:  &exit,
			Lots:  1.0,
		}

		CalculateTradeMetrics(trade)

		// Assert R:R is not set (empty string)
		if trade.RR != "" {
			t.Errorf("expected RR to be empty without stop loss, got %s", trade.RR)
		}

		// But pips and P/L should still be calculated
		if trade.Pips == nil {
			t.Error("expected pips to be calculated")
		}
		if trade.PL == nil {
			t.Error("expected PL to be calculated")
		}
	})
}

// Helper function to create float pointer
func floatPtr(f float64) *float64 {
	return &f
}
