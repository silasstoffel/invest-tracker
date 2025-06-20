package config

import (
	"fmt"
	"os"
	"strings"
)

type CloudflareConfig struct {
	AccountId           string
	InvestmentTrackDbId string
	ApiKeyParamName     string // Identifier to read API Key value
}

type Config struct {
	Env                      string
	CreateInvestmentQueueURL string
	Cloudflare               *CloudflareConfig
}

func NewConfigFromEnvVars() *Config {
	e := strings.ToLower(os.Getenv("ENVIRONMENT"))
	env := "development"
	getEnvPrefix := "DEV"
	apiKeyParamName := "/invest-track-dev/cloudflare/api-key"

	if e == "prod" || e == "production" || e == "prd" {
		env = "production"
		getEnvPrefix = "PROD"
		apiKeyParamName = "/invest-track-prod/cloudflare/api-key"
	}

	return &Config{
		Env:                      env,
		CreateInvestmentQueueURL: os.Getenv("CREATE_INVESTMENT_QUEUE_URL"),
		Cloudflare: &CloudflareConfig{
			AccountId:           os.Getenv(fmt.Sprintf("CLOUDFLARE_ACCOUNT_ID_%s", getEnvPrefix)),
			InvestmentTrackDbId: os.Getenv(fmt.Sprintf("CLOUDFLARE_DB_ID_PROD_%s", getEnvPrefix)),
			ApiKeyParamName:     apiKeyParamName,
		},
	}
}
