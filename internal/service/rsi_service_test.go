package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/deba0208/stock-rsi-dashboard/internal/models"
	"github.com/deba0208/stock-rsi-dashboard/internal/service"
)

// --- Stub provider helpers ---

type stubProvider struct {
	daily   []models.Candle
	weekly  []models.Candle
	monthly []models.Candle
	err     error
}

func (s *stubProvider) GetDailyCandles(_ string) ([]models.Candle, error) {
	return s.daily, s.err
}
func (s *stubProvider) GetWeeklyCandles(_ string) ([]models.Candle, error) {
	return s.weekly, s.err
}
func (s *stubProvider) GetMonthlyCandles(_ string) ([]models.Candle, error) {
	return s.monthly, s.err
}

// 15 rising candles (enough for 14-period RSI)
func risingCandles() []models.Candle {
	prices := []float64{
		100, 101, 102, 103, 104, 105,
		106, 107, 108, 109, 110, 111,
		112, 113, 114,
	}
	candles := make([]models.Candle, len(prices))
	for i, p := range prices {
		candles[i] = models.Candle{Date: time.Now(), Close: p}
	}
	return candles
}

// --- Tests ---

func TestRSIService_Daily_Success(t *testing.T) {
	provider := &stubProvider{daily: risingCandles()}
	svc := service.NewRSIService(provider)

	result, err := svc.RSI("RELIANCE", service.Daily)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result <= 0 || result > 100 {
		t.Errorf("RSI out of range, got %.2f", result)
	}
}

func TestRSIService_Weekly_Success(t *testing.T) {
	provider := &stubProvider{weekly: risingCandles()}
	svc := service.NewRSIService(provider)

	result, err := svc.RSI("TCS", service.Weekly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result <= 0 || result > 100 {
		t.Errorf("RSI out of range, got %.2f", result)
	}
}

func TestRSIService_Monthly_Success(t *testing.T) {
	provider := &stubProvider{monthly: risingCandles()}
	svc := service.NewRSIService(provider)

	result, err := svc.RSI("INFY", service.Monthly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result <= 0 || result > 100 {
		t.Errorf("RSI out of range, got %.2f", result)
	}
}

func TestRSIService_InvalidTimeframe_ReturnsError(t *testing.T) {
	provider := &stubProvider{daily: risingCandles()}
	svc := service.NewRSIService(provider)

	_, err := svc.RSI("RELIANCE", "invalid")
	if err == nil {
		t.Error("expected error for invalid timeframe, got nil")
	}
}

func TestRSIService_ProviderError_ReturnsError(t *testing.T) {
	provider := &stubProvider{err: errors.New("provider down")}
	svc := service.NewRSIService(provider)

	_, err := svc.RSI("RELIANCE", service.Daily)
	if err == nil {
		t.Error("expected error when provider fails, got nil")
	}
}

func TestRSIService_EmptyCandles_ReturnsError(t *testing.T) {
	provider := &stubProvider{daily: []models.Candle{}} // empty
	svc := service.NewRSIService(provider)

	_, err := svc.RSI("RELIANCE", service.Daily)
	if err == nil {
		t.Error("expected error for empty candles, got nil")
	}
}

func TestRSIService_AllGains_Returns100(t *testing.T) {
	provider := &stubProvider{daily: risingCandles()}
	svc := service.NewRSIService(provider)

	result, err := svc.RSI("RELIANCE", service.Daily)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 100 {
		t.Errorf("expected 100 for all-gains candles, got %.2f", result)
	}
}
