package analytics

// Analytics represents the trading analytics/metrics for a user
type Analytics struct {
	// Performance Metrics
	TotalPL        float64 // Total Profit/Loss
	WinRate        float64 // Win rate percentage
	TotalTrades    int64   // Total number of trades
	WinningTrades  int64   // Number of winning trades
	LosingTrades   int64   // Number of losing trades
	AvgWin         float64 // Average winning trade
	AvgLoss        float64 // Average losing trade
	ProfitFactor   float64 // Profit Factor (gross profit / gross loss)
	SharpeRatio    float64 // Sharpe Ratio
	MaxDrawdown    float64 // Maximum drawdown
	LargestWin     float64 // Largest winning trade
	LargestLoss    float64 // Largest losing trade

	// Additional Metrics
	AvgRR          float64 // Average Risk:Reward ratio
	ConsecutiveWins int64   // Current consecutive wins
	ConsecutiveLosses int64 // Current consecutive losses
	BestStreak     int64   // Best winning streak
	WorstStreak    int64   // Worst losing streak
}
