package persistence

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/raihanstark/trade-journal/internal/db"
	infradb "github.com/raihanstark/trade-journal/internal/db"
	"github.com/raihanstark/trade-journal/internal/domain/trade"
)

type TradeRepository struct {
	queries *db.Queries
}

func NewTradeRepository(queries *db.Queries) *TradeRepository {
	return &TradeRepository{queries: queries}
}

func (r *TradeRepository) GetByAccountID(ctx context.Context, accountID int64, userID int64) ([]*trade.Trade, error) {
	results, err := r.queries.GetTradesByAccountID(ctx, db.GetTradesByAccountIDParams{
		AccountID: sql.NullInt32{Int32: int32(accountID), Valid: true},
		UserID:    int32(userID),
	})
	if err != nil {
		return nil, err
	}

	trades := make([]*trade.Trade, len(results))
	for i, result := range results {
		strategies, err := r.queries.GetTradeStrategies(ctx, result.ID)
		if err != nil {
			return nil, err
		}
		trades[i] = r.toDomain(&result, strategies)
	}

	return trades, nil
}

func (r *TradeRepository) Create(ctx context.Context, t *trade.Trade) (*trade.Trade, error) {
	result, err := r.queries.CreateTrade(ctx, db.CreateTradeParams{
		UserID:     int32(t.UserID),
		AccountID:  int32ToNullInt32(t.AccountID),
		Date:       t.Date,
		Time:       t.Time,
		Pair:       infradb.StringToNullString(t.Pair),
		Type:       db.TradeType(t.Type),
		Entry:      floatToNullString(t.Entry),
		Exit:       floatPtrToNullString(t.Exit),
		Lots:       floatToNullString(t.Lots),
		Pips:       floatPtrToNullString(t.Pips),
		Pl:         floatPtrToNullString(t.PL),
		Rr:         infradb.StringToNullString(t.RR),
		Status:     db.TradeStatus(t.Status),
		StopLoss:   floatPtrToNullString(t.StopLoss),
		TakeProfit: floatPtrToNullString(t.TakeProfit),
		Notes:      infradb.StringToNullString(t.Notes),
		Mistakes:   infradb.StringToNullString(t.Mistakes),
		Amount:     floatPtrToNullString(t.Amount),
	})
	if err != nil {
		return nil, err
	}

	// Add strategies
	for _, strategy := range t.Strategies {
		err := r.queries.AddTradeStrategy(ctx, db.AddTradeStrategyParams{
			TradeID:    result.ID,
			StrategyID: int32(strategy.ID),
		})
		if err != nil {
			return nil, err
		}
	}

	// Fetch strategies for the created trade
	strategies, err := r.queries.GetTradeStrategies(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	return r.toDomain(&result, strategies), nil
}

func (r *TradeRepository) GetByID(ctx context.Context, id int64, userID int64) (*trade.Trade, error) {
	result, err := r.queries.GetTradeByID(ctx, db.GetTradeByIDParams{
		ID:     int32(id),
		UserID: int32(userID),
	})
	if err != nil {
		return nil, err
	}

	strategies, err := r.queries.GetTradeStrategies(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	return r.toDomain(&result, strategies), nil
}

func (r *TradeRepository) GetByUserID(ctx context.Context, userID int64) ([]*trade.Trade, error) {
	results, err := r.queries.GetTradesByUserID(ctx, int32(userID))
	if err != nil {
		return nil, err
	}

	trades := make([]*trade.Trade, len(results))
	for i, result := range results {
		strategies, err := r.queries.GetTradeStrategies(ctx, result.ID)
		if err != nil {
			return nil, err
		}
		trades[i] = r.toDomain(&result, strategies)
	}

	return trades, nil
}

func (r *TradeRepository) GetByUserIDAndDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*trade.Trade, error) {
	results, err := r.queries.GetTradesByUserIDAndDateRange(ctx, db.GetTradesByUserIDAndDateRangeParams{
		UserID:   int32(userID),
		Date:     startDate,
		Date_2:   endDate,
	})
	if err != nil {
		return nil, err
	}

	trades := make([]*trade.Trade, len(results))
	for i, result := range results {
		strategies, err := r.queries.GetTradeStrategies(ctx, result.ID)
		if err != nil {
			return nil, err
		}
		trades[i] = r.toDomain(&result, strategies)
	}

	return trades, nil
}

func (r *TradeRepository) GetByAccountIDAndDateRange(ctx context.Context, accountID int64, userID int64, startDate, endDate time.Time) ([]*trade.Trade, error) {
	results, err := r.queries.GetTradesByAccountIDAndDateRange(ctx, db.GetTradesByAccountIDAndDateRangeParams{
		AccountID: sql.NullInt32{Int32: int32(accountID), Valid: true},
		UserID:    int32(userID),
		Date:      startDate,
		Date_2:    endDate,
	})
	if err != nil {
		return nil, err
	}

	trades := make([]*trade.Trade, len(results))
	for i, result := range results {
		strategies, err := r.queries.GetTradeStrategies(ctx, result.ID)
		if err != nil {
			return nil, err
		}
		trades[i] = r.toDomain(&result, strategies)
	}

	return trades, nil
}

func (r *TradeRepository) Update(ctx context.Context, t *trade.Trade) (*trade.Trade, error) {
	result, err := r.queries.UpdateTrade(ctx, db.UpdateTradeParams{
		ID:         int32(t.ID),
		AccountID:  int32ToNullInt32(t.AccountID),
		Date:       t.Date,
		Time:       t.Time,
		Pair:       infradb.StringToNullString(t.Pair),
		Type:       db.TradeType(t.Type),
		Entry:      floatToNullString(t.Entry),
		Exit:       floatPtrToNullString(t.Exit),
		Lots:       floatToNullString(t.Lots),
		Pips:       floatPtrToNullString(t.Pips),
		Pl:         floatPtrToNullString(t.PL),
		Rr:         infradb.StringToNullString(t.RR),
		Status:     db.TradeStatus(t.Status),
		StopLoss:   floatPtrToNullString(t.StopLoss),
		TakeProfit: floatPtrToNullString(t.TakeProfit),
		Notes:      infradb.StringToNullString(t.Notes),
		Mistakes:   infradb.StringToNullString(t.Mistakes),
		Amount:     floatPtrToNullString(t.Amount),
		UserID:     int32(t.UserID),
	})
	if err != nil {
		return nil, err
	}

	// Delete existing strategies and add new ones
	err = r.queries.DeleteTradeStrategies(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	for _, strategy := range t.Strategies {
		err := r.queries.AddTradeStrategy(ctx, db.AddTradeStrategyParams{
			TradeID:    result.ID,
			StrategyID: int32(strategy.ID),
		})
		if err != nil {
			return nil, err
		}
	}

	strategies, err := r.queries.GetTradeStrategies(ctx, result.ID)
	if err != nil {
		return nil, err
	}

	return r.toDomain(&result, strategies), nil
}

func (r *TradeRepository) Delete(ctx context.Context, id int64, userID int64) error {
	return r.queries.DeleteTrade(ctx, db.DeleteTradeParams{
		ID:     int32(id),
		UserID: int32(userID),
	})
}

func (r *TradeRepository) toDomain(t *db.Trade, strategies []db.Strategy) *trade.Trade {
	domainStrategies := make([]trade.Strategy, len(strategies))
	for i, s := range strategies {
		domainStrategies[i] = trade.Strategy{
			ID:          int64(s.ID),
			Name:        s.Name,
			Description: infradb.NullStringToString(s.Description),
		}
	}

	return &trade.Trade{
		ID:         int64(t.ID),
		UserID:     int64(t.UserID),
		AccountID:  nullInt32ToInt64Ptr(t.AccountID),
		Date:       t.Date,
		Time:       t.Time,
		Pair:       infradb.NullStringToString(t.Pair),
		Type:       trade.TradeType(t.Type),
		Entry:      nullStringToFloat(t.Entry),
		Exit:       nullStringToFloatPtr(t.Exit),
		Lots:       nullStringToFloat(t.Lots),
		Pips:       nullStringToFloatPtr(t.Pips),
		PL:         nullStringToFloatPtr(t.Pl),
		RR:         infradb.NullStringToString(t.Rr),
		Status:     trade.TradeStatus(t.Status),
		StopLoss:   nullStringToFloatPtr(t.StopLoss),
		TakeProfit: nullStringToFloatPtr(t.TakeProfit),
		Notes:      infradb.NullStringToString(t.Notes),
		Mistakes:   infradb.NullStringToString(t.Mistakes),
		Amount:     nullStringToFloatPtr(t.Amount),
		Strategies: domainStrategies,
		CreatedAt:  t.CreatedAt.Time,
		UpdatedAt:  t.UpdatedAt.Time,
	}
}

// Helper functions for type conversion
func int32ToNullInt32(i *int64) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*i), Valid: true}
}

func nullInt32ToInt64Ptr(n sql.NullInt32) *int64 {
	if !n.Valid {
		return nil
	}
	v := int64(n.Int32)
	return &v
}

func floatToNullString(f float64) sql.NullString {
	if f == 0 {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: formatFloat(f), Valid: true}
}

func floatPtrToNullString(f *float64) sql.NullString {
	if f == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: formatFloat(*f), Valid: true}
}

func nullStringToFloat(n sql.NullString) float64 {
	if !n.Valid {
		return 0
	}
	return parseFloat(n.String)
}

func nullStringToFloatPtr(n sql.NullString) *float64 {
	if !n.Valid {
		return nil
	}
	v := parseFloat(n.String)
	return &v
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
