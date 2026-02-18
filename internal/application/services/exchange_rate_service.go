package services

import (
	"context"
	"fmt"
	"payment-system/internal/domain/entities"
	"payment-system/internal/domain/repositories"
	"time"

	"github.com/shopspring/decimal"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         decimal.Decimal
	EffectiveAt  time.Time
	Source       string
}

type ExchangeRateService interface {
	GetExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (*ExchangeRate, error)
	ConvertAmount(ctx context.Context, amount decimal.Decimal, fromCurrency, toCurrency string) (*decimal.Decimal, error)
}

type exchangeRateService struct {
	ledgerRepo repositories.LedgerRepository
	// In a real implementation, you'd inject an external exchange rate provider
	// For now, we'll use a simple mock implementation
}

func NewExchangeRateService(ledgerRepo repositories.LedgerRepository) ExchangeRateService {
	return &exchangeRateService{
		ledgerRepo: ledgerRepo,
	}
}

func (s *exchangeRateService) GetExchangeRate(ctx context.Context, fromCurrency, toCurrency string) (*ExchangeRate, error) {
	if fromCurrency == toCurrency {
		return &ExchangeRate{
			FromCurrency: fromCurrency,
			ToCurrency:   toCurrency,
			Rate:         decimal.NewFromInt(1),
			EffectiveAt:  time.Now(),
			Source:       "internal",
		}, nil
	}

	// TODO: Integrate with real exchange rate API (e.g., Fixer.io, ExchangeRate-API)
	// For now, return mock rates
	rate := s.getMockRate(fromCurrency, toCurrency)
	if rate == nil {
		return nil, fmt.Errorf("exchange rate not available for %s to %s", fromCurrency, toCurrency)
	}

	return &ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         *rate,
		EffectiveAt:  time.Now(),
		Source:       "mock",
	}, nil
}

func (s *exchangeRateService) ConvertAmount(ctx context.Context, amount decimal.Decimal, fromCurrency, toCurrency string) (*decimal.Decimal, error) {
	rate, err := s.GetExchangeRate(ctx, fromCurrency, toCurrency)
	if err != nil {
		return nil, err
	}

	converted := amount.Mul(rate.Rate)
	return &converted, nil
}

// Mock exchange rates - replace with real API integration
func (s *exchangeRateService) getMockRate(from, to string) *decimal.Decimal {
	rates := map[string]decimal.Decimal{
		"USD_EUR": decimal.NewFromFloat(0.92),
		"EUR_USD": decimal.NewFromFloat(1.09),
		"USD_GBP": decimal.NewFromFloat(0.79),
		"GBP_USD": decimal.NewFromFloat(1.27),
		"EUR_GBP": decimal.NewFromFloat(0.86),
		"GBP_EUR": decimal.NewFromFloat(1.16),
	}

	key := from + "_" + to
	if rate, ok := rates[key]; ok {
		return &rate
	}
	return nil
}
