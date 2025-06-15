package investment

import "time"

const (
	// investiment types
	FiiInvestmentType   = "FII"
	StockInvestmentType = "STOCK"
	ReitsInvestmentType = "REITS"
	BondInvestmentType  = "BOND"

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
	DueDate        int       `json:"dueDate"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CreateInvestmentInput struct {
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
}

type CreateInvestmentOutput struct{}
