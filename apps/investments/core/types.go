package investment_core

import (
	"time"
)

const (
	// investiment types
	FiiInvestmentType   = "fii"
	StockInvestmentType = "stock"
	ReitInvestmentType  = "reit"
	BondInvestmentType  = "bond"

	// operation types
	BuyOperationType  = "buy"
	SellOperationType = "sell"

	BondIndexCDI    = "cdi"
	BondIndexIPCA   = "ipca"
	BondIndexSELIC  = "selic"
	BondIndexPrefix = "prefixed"
)

type InvestmentEntity struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Symbol         string    `json:"symbol"`
	BondIndex      string    `json:"bondIndex,omitempty"`
	BondRate       float64   `json:"bondRate,omitempty"`
	Quantity       int       `json:"quantity"`
	UnitPrice      float64   `json:"unitPrice"`
	TotalValue     float64   `json:"totalValue"`
	Cost           float64   `json:"cost"`
	OperationType  string    `json:"operationType"`
	OperationDate  string    `json:"operationDate"`
	OperationYear  int       `json:"operationYear"`
	OperationMonth int       `json:"operationMonth"`
	DueDate        string    `json:"dueDate"`
	Brokerage      string    `json:"brokerage"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CreateInvestmentInput struct {
	Type          string  `json:"type"`
	Symbol        string  `json:"symbol"`
	BondIndex     string  `json:"bondIndex,omitempty"`
	BondRate      float64 `json:"bondRate,omitempty"`
	Quantity      int     `json:"quantity"`
	TotalValue    float64 `json:"totalValue"`
	Cost          float64 `json:"cost"`
	OperationType string  `json:"operationType"`
	OperationDate string  `json:"operationDate"`
	DueDate       string  `json:"dueDate"`
	Brokerage     string  `json:"brokerage"`
}

type CreateInvestmentOutput struct {
	Message string `json:"message"`
}
