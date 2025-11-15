package trade

import "time"

type TradeType string

const (
	TradeTypeBuy      TradeType = "BUY"
	TradeTypeSell     TradeType = "SELL"
	TradeTypeDeposit  TradeType = "DEPOSIT"
	TradeTypeWithdraw TradeType = "WITHDRAW"
)

type TradeStatus string

const (
	TradeStatusOpen   TradeStatus = "open"
	TradeStatusClosed TradeStatus = "closed"
)

type Trade struct {
	ID         int64
	UserID     int64
	AccountID  *int64
	Date       time.Time
	Time       time.Time
	Pair       string
	Type       TradeType
	Entry      float64
	Exit       *float64
	Lots       float64
	Pips       *float64
	PL         *float64
	RR         *float64
	Status     TradeStatus
	StopLoss   *float64
	TakeProfit *float64
	Notes      string
	Mistakes   string
	Amount     *float64
	Strategies []Strategy
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Strategy struct {
	ID          int64
	Name        string
	Description string
}
