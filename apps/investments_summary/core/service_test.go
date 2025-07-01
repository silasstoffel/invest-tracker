package investment_summary_core

import (
	"testing"
)

func TestCalculateAverageCost(t *testing.T) {
	investments := []InvestmentCreatedInput{
		{
			OperationType: "buy",
			Quantity:      10,
			TotalValue:    100.0,
		},
		{
			OperationType: "buy",
			Quantity:      15,
			TotalValue:    160.0,
		},
		{
			OperationType: "sell",
			Quantity:      5,
			TotalValue:    55.0,
		},
	}

	created := InvestmentCreatedInput{
		Quantity:   5,
		TotalValue: 55.0,
	}

	result := CalculateAverageCost(created, investments)

	expectedQuantity := 20
	expectedTotalValue := 208.0
	expectedAveragePrice := 10.4

	if result.Quantity != expectedQuantity {
		t.Errorf("expected quantity %d, got %d", expectedQuantity, result.Quantity)
	}
	if result.TotalValue != expectedTotalValue {
		t.Errorf("expected total value %.2f, got %.2f", expectedTotalValue, result.TotalValue)
	}
	if result.AveragePrice != expectedAveragePrice {
		t.Errorf("expected average price %.2f, got %.2f", expectedAveragePrice, result.AveragePrice)
	}
}
