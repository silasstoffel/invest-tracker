package investment_summary_core

import "log"

func CalculateAverageCost(createdInvestment InvestmentCreatedInput, investments []InvestmentCreatedInput) CalculateAverageCostOutput {
	if len(investments) == 0 {
		log.Print("It is the first operation to summarize")
		return CalculateAverageCostOutput{
			Quantity:     createdInvestment.Quantity,
			TotalValue:   createdInvestment.TotalValue,
			AveragePrice: createdInvestment.TotalValue / float64(createdInvestment.Quantity),
		}
	}
	return handleAverageCost(investments)
}

func handleAverageCost(investments []InvestmentCreatedInput) CalculateAverageCostOutput {
	quantity := 0
	var totalValue float64 = 0
	var averagePrice float64 = 0

	for _, investment := range investments {
		if investment.OperationType == "buy" {
			quantity += investment.Quantity
			totalValue += investment.TotalValue
			averagePrice = totalValue / float64(quantity)
		} else {
			costToRemove := float64(investment.Quantity) * averagePrice
			totalValue -= costToRemove
			quantity -= investment.Quantity

			if quantity > 0 {
				averagePrice = totalValue / float64(quantity)
			} else {
				averagePrice = 0
			}
		}
	}

	return CalculateAverageCostOutput{
		Quantity:     quantity,
		TotalValue:   totalValue,
		AveragePrice: averagePrice,
	}
}
