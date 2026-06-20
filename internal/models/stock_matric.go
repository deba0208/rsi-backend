package models

type StockMetric struct {
	Symbol     string  `json:"symbol"`
	Price      float64 `json:"price"`
	DailyRSI   float64 `json:"dailyRsi"`
	WeeklyRSI  float64 `json:"weeklyRsi"`
	MonthlyRSI float64 `json:"monthlyRsi"`
}
