package trade

import "time"

type TradeDTO struct {
	ID         int64      `json:"id"`
	AccountID  *int64     `json:"account_id"`
	Date       string     `json:"date"`
	Time       string     `json:"time"`
	Pair       string     `json:"pair"`
	Type       string     `json:"type"`
	Entry      float64    `json:"entry"`
	Exit       *float64   `json:"exit"`
	Lots       float64    `json:"lots"`
	Pips       *float64   `json:"pips"`
	PL         *float64   `json:"pl"`
	RR         *float64   `json:"rr"`
	Status     string     `json:"status"`
	StopLoss   *float64   `json:"stop_loss"`
	TakeProfit *float64   `json:"take_profit"`
	Notes      string     `json:"notes"`
	Mistakes   string     `json:"mistakes"`
	Amount     *float64   `json:"amount"`
	Strategies []Strategy `json:"strategies"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Strategy struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CreateTradeRequest struct {
	AccountID   *int64   `json:"account_id"`
	Date        string   `json:"date"`
	Time        string   `json:"time"`
	Pair        string   `json:"pair"`
	Type        string   `json:"type"`
	Entry       float64  `json:"entry"`
	Exit        *float64 `json:"exit"`
	Lots        float64  `json:"lots"`
	Pips        *float64 `json:"pips"`
	PL          *float64 `json:"pl"`
	RR          *float64 `json:"rr"`
	Status      string   `json:"status"`
	StopLoss    *float64 `json:"stop_loss"`
	TakeProfit  *float64 `json:"take_profit"`
	Notes       string   `json:"notes"`
	Mistakes    string   `json:"mistakes"`
	Amount      *float64 `json:"amount"`
	StrategyIDs []int64  `json:"strategy_ids"`
}

type UpdateTradeRequest struct {
	AccountID   *int64   `json:"account_id"`
	Date        string   `json:"date"`
	Time        string   `json:"time"`
	Pair        string   `json:"pair"`
	Type        string   `json:"type"`
	Entry       float64  `json:"entry"`
	Exit        *float64 `json:"exit"`
	Lots        float64  `json:"lots"`
	Pips        *float64 `json:"pips"`
	PL          *float64 `json:"pl"`
	RR          *float64 `json:"rr"`
	Status      string   `json:"status"`
	StopLoss    *float64 `json:"stop_loss"`
	TakeProfit  *float64 `json:"take_profit"`
	Notes       string   `json:"notes"`
	Mistakes    string   `json:"mistakes"`
	Amount      *float64 `json:"amount"`
	StrategyIDs []int64  `json:"strategy_ids"`
}
