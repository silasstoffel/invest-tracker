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
	EtfInvestmentType   = "etf"

	// operation types
	BuyOperationType  = "buy"
	SellOperationType = "sell"

	BondIndexCDI    = "cdi"
	BondIndexIPCA   = "ipca"
	BondIndexSELIC  = "selic"
	BondIndexPrefix = "prefixed"

	HybridRedemption     = "hybrid"
	AnyTimeRedemption    = "any_time"
	AtMaturityRedemption = "at_maturity"
)

type InvestmentEntity struct {
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

type CreateInvestmentInput struct {
	Type                 string  `json:"type"`
	Symbol               string  `json:"symbol"`
	BondIndex            string  `json:"bondIndex,omitempty"`
	BondRate             float64 `json:"bondRate,omitempty"`
	Quantity             float64 `json:"quantity"`
	TotalValue           float64 `json:"totalValue"`
	Cost                 float64 `json:"cost"`
	OperationType        string  `json:"operationType"`
	OperationDate        string  `json:"operationDate"`
	DueDate              string  `json:"dueDate"`
	Brokerage            string  `json:"brokerage"`
	Note                 string  `json:"note"`
	RedemptionPolicyType string  `json:"redemptionPolicyType"`
	// required for bond investment and sell operation type
	SellInvestmentId string `json:"sellInvestmentId,omitempty"`
}

type CreateInvestmentOutput struct {
	Message string `json:"message"`
}
