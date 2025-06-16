package investment_core

import "time"

const (
	// investiment types
	FiiInvestmentType   = "fii"
	StockInvestmentType = "stock"
	ReitInvestmentType  = "reit"
	BondInvestmentType  = "bond"

	// operation types
	BuyOperationType  = "buy"
	SellOperationType = "sell"
)

type InvestmentEntity struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Symbol         string    `json:"symbol"`
	Quantity       int       `json:"quantity"`
	UnitPrice      float64   `json:"unitPrice"`
	TotalValue     float64   `json:"totalValue"`
	Cost           float64   `json:"cost"`
	OperationType  string    `json:"operationType"`
	OperationDate  time.Time `json:"operationDate"`
	OperationYear  int       `json:"operationYear"`
	OperationMonth int       `json:"operationMonth"`
	DueDate        time.Time `json:"dueDate"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CreateInvestmentInput struct {
	Type          string    `json:"type"`
	Symbol        string    `json:"symbol"`
	Quantity      int       `json:"quantity"`
	TotalValue    float64   `json:"totalValue"`
	Cost          float64   `json:"cost"`
	OperationType string    `json:"operationType"`
	OperationDate time.Time `json:"operationDate"`
	DueDate       time.Time `json:"dueDate"`
}

type CreateInvestmentOutput struct {
	Message string `json:"message"`
}
