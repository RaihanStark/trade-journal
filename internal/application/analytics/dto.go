package analytics

type AnalyticsDTO struct {
	TotalPL           float64 `json:"total_pl"`
	WinRate           float64 `json:"win_rate"`
	TotalTrades       int64   `json:"total_trades"`
	WinningTrades     int64   `json:"winning_trades"`
	LosingTrades      int64   `json:"losing_trades"`
	AvgWin            float64 `json:"avg_win"`
	AvgLoss           float64 `json:"avg_loss"`
	ProfitFactor      float64 `json:"profit_factor"`
	SharpeRatio       float64 `json:"sharpe_ratio"`
	MaxDrawdown       float64 `json:"max_drawdown"`
	LargestWin        float64 `json:"largest_win"`
	LargestLoss       float64 `json:"largest_loss"`
	AvgRR             float64 `json:"avg_rr"`
	ConsecutiveWins   int64   `json:"consecutive_wins"`
	ConsecutiveLosses int64   `json:"consecutive_losses"`
	BestStreak        int64   `json:"best_streak"`
	WorstStreak       int64   `json:"worst_streak"`
}
