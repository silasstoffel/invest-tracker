package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type CloudflareConfig struct {
	AccountId           string
	InvestmentTrackDbId string
	ApiKeyParamName     string // Identifier to read API Key value
}

type Aws struct{}

type TelegramConfig struct {
	Token  string
	ChatId string
}

type Config struct {
	Env                           string
	CreateInvestmentQueueURL      string
	CalculateAveragePriceQueueURL string
	Cloudflare                    CloudflareConfig
	Aws                           *Aws
	TelegramConfig                *TelegramConfig
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
		Env:                           env,
		CreateInvestmentQueueURL:      os.Getenv("CREATE_INVESTMENT_QUEUE_URL"),
		CalculateAveragePriceQueueURL: os.Getenv("CALCULATE_AVERAGE_PRICE_QUEUE_URL"),
		Cloudflare: CloudflareConfig{
			AccountId:           os.Getenv(fmt.Sprintf("CLOUDFLARE_ACCOUNT_ID_%s", getEnvPrefix)),
			InvestmentTrackDbId: os.Getenv(fmt.Sprintf("CLOUDFLARE_DB_ID_%s", getEnvPrefix)),
			ApiKeyParamName:     apiKeyParamName,
		},
		Aws: &Aws{},
		TelegramConfig: &TelegramConfig{
			Token:  os.Getenv("TELEGRAM_TOKEN"),
			ChatId: os.Getenv("TELEGRAM_CHAT_ID"),
		},
	}
}

func (a *Aws) LoadDefaultConfig() (aws.Config, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, fmt.Errorf("%s", fmt.Sprintf("Failure to load aws config: %v", err))
	}
	return cfg, nil
}
