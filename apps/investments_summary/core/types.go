package investment_summary_core

import (
	"time"
)

type InvestmentCreatedInput struct {
	ID                   string    `json:"id"`
	Type                 string    `json:"type"`
	Symbol               string    `json:"symbol"`
	BondIndex            string    `json:"bondIndex,omitempty"`
	BondRate             float64   `json:"bondRate,omitempty"`
	Quantity             float64   `json:"quantity"`
	UnitPrice            float64   `json:"unitPrice"`
	TotalValue           float64   `json:"totalValue"`
	Cost                 float64   `json:"cost"`
	OperationType        string    `json:"operationType"`
	OperationDate        string    `json:"operationDate"`
	OperationYear        int       `json:"operationYear"`
	OperationMonth       int       `json:"operationMonth"`
	DueDate              string    `json:"dueDate"`
	Brokerage            string    `json:"brokerage"`
	Note                 string    `json:"note"`
	RedemptionPolicyType string    `json:"redemptionPolicyType"`
	SellInvestmentId     string    `json:"sellInvestmentId,omitempty"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type InvestmentSummaryEntity struct {
	ID                   string    `json:"id"`
	Type                 string    `json:"type"`
	Symbol               string    `json:"symbol"`
	BondIndex            string    `json:"bondIndex,omitempty"`
	BondRate             float64   `json:"bondRate,omitempty"`
	Quantity             float64   `json:"quantity"`
	AveragePrice         float64   `json:"averagePrice"`
	TotalValue           float64   `json:"totalValue"`
	Cost                 float64   `json:"cost"`
	LastTransactionDate  time.Time `json:"lastOperationDate"`
	LastTransactionType  string    `json:"lastOperationType"`
	DueDate              string    `json:"dueDate"`
	Brokerage            string    `json:"brokerage"`
	RedemptionPolicyType string    `json:"redemptionPolicyType"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type CalculateAverageCostOutput struct {
	Quantity     float64
	TotalValue   float64
	AveragePrice float64
}
