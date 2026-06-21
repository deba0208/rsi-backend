package repository

import (
	"context"
	"encoding/json"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/redis/go-redis/v9"
)

type MetricRepository struct {
	client *redis.Client
}

func NewMetricRepository(client *redis.Client) *MetricRepository {
	return &MetricRepository{client: client}
}

// SaveMetric stores metric as JSON and updates RSI sorted set rankings
func (r *MetricRepository) SaveMetric(metric models.StockMetric) error {
	ctx := context.Background()

	// Store full metric as JSON: "metric:RELIANCE" → {...}
	data, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	if err := r.client.Set(ctx, "metric:"+metric.Symbol, data, 0).Err(); err != nil {
		return err
	}

	// Update RSI sorted set rankings
	return r.SaveRanking(
		metric.Symbol,
		metric.DailyRSI,
		metric.WeeklyRSI,
		metric.MonthlyRSI,
	)
}

// SaveRanking adds RSI scores to sorted sets rsi:daily, rsi:weekly, rsi:monthly
func (r *MetricRepository) SaveRanking(
	symbol string,
	daily float64,
	weekly float64,
	monthly float64,
) error {

	ctx := context.Background()

	if err := r.client.ZAdd(ctx, "rsi:daily", redis.Z{
		Score:  daily,
		Member: symbol,
	}).Err(); err != nil {
		return err
	}

	if err := r.client.ZAdd(ctx, "rsi:weekly", redis.Z{
		Score:  weekly,
		Member: symbol,
	}).Err(); err != nil {
		return err
	}

	if err := r.client.ZAdd(ctx, "rsi:monthly", redis.Z{
		Score:  monthly,
		Member: symbol,
	}).Err(); err != nil {
		return err
	}

	return nil
}

// GetMetric retrieves a stock metric by symbol
func (r *MetricRepository) GetMetric(symbol string) (*models.StockMetric, error) {
	ctx := context.Background()

	data, err := r.client.Get(ctx, "metric:"+symbol).Bytes()
	if err != nil {
		return nil, err
	}

	var metric models.StockMetric
	if err := json.Unmarshal(data, &metric); err != nil {
		return nil, err
	}

	return &metric, nil
}

func (r *MetricRepository) GetTop50ByCriteria(criteria string) ([]string, error) {

	ctx := context.Background()

	// ZRevRange returns members in descending score order (highest RSI first)
	return r.client.ZRevRange(
		ctx,
		"rsi:"+criteria,
		0,
		49,
	).Result()
}
