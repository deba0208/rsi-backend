package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/redis/go-redis/v9"
)

type MetricRepository struct {
	client *redis.Client
}

func NewMetricRepository(client *redis.Client) *MetricRepository {
	return &MetricRepository{client: client}
}

// SaveMetric stores metric as a Redis Hash and updates RSI sorted set rankings
func (r *MetricRepository) SaveMetric(metric models.StockMetric) error {
	ctx := context.Background()

	// Store each field as a named hash field
	if err := r.client.HSet(ctx, "metric:"+metric.Symbol,
		"symbol", metric.Symbol,
		"price", fmt.Sprintf("%f", metric.Price),
		"dailyRsi", fmt.Sprintf("%f", metric.DailyRSI),
		"weeklyRsi", fmt.Sprintf("%f", metric.WeeklyRSI),
		"monthlyRsi", fmt.Sprintf("%f", metric.MonthlyRSI),
	).Err(); err != nil {
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

	// HGetAll returns map[string]string of all hash fields
	fields, err := r.client.HGetAll(ctx, "metric:"+symbol).Result()
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("metric not found for symbol: %s", symbol)
	}

	price, _ := strconv.ParseFloat(fields["price"], 64)
	dailyRSI, _ := strconv.ParseFloat(fields["dailyRsi"], 64)
	weeklyRSI, _ := strconv.ParseFloat(fields["weeklyRsi"], 64)
	monthlyRSI, _ := strconv.ParseFloat(fields["monthlyRsi"], 64)

	return &models.StockMetric{
		Symbol:     fields["symbol"],
		Price:      price,
		DailyRSI:   dailyRSI,
		WeeklyRSI:  weeklyRSI,
		MonthlyRSI: monthlyRSI,
	}, nil
}

func (r *MetricRepository) GetTopByTimeFrame(
	timeFrame string,
	limit int64,
) ([]models.StockMetric, error) {
	rankingKey := fmt.Sprintf(
		"rsi:%s",
		timeFrame,
	)

	symbols, err :=
		r.client.ZRange(
			context.Background(),
			rankingKey,
			0,
			limit-1,
		).Result()

	if err != nil {
		return nil, err
	}

	metrics :=
		make(
			[]models.StockMetric,
			0,
			len(symbols),
		)

	for _, symbol := range symbols {
		metricKey :=
			fmt.Sprintf(
				"metric:%s",
				symbol,
			)

		values, err :=
			r.client.HGetAll(
				context.Background(),
				metricKey,
			).Result()

		if err != nil {
			continue
		}

		price, _ :=
			strconv.ParseFloat(
				values["price"],
				64,
			)

		dailyRsi, _ :=
			strconv.ParseFloat(
				values["dailyRsi"],
				64,
			)

		weeklyRsi, _ :=
			strconv.ParseFloat(
				values["weeklyRsi"],
				64,
			)

		monthlyRsi, _ :=
			strconv.ParseFloat(
				values["monthlyRsi"],
				64,
			)

		metric := models.StockMetric{
			Symbol:     values["symbol"],
			Price:      price,
			DailyRSI:   dailyRsi,
			WeeklyRSI:  weeklyRsi,
			MonthlyRSI: monthlyRsi,
		}
		metrics =
			append(
				metrics,
				metric,
			)
	}

	return metrics, nil
}
