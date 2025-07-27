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
	ApiKey              string
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

	if e == "prod" || e == "production" || e == "prd" {
		env = "production"
	}

	return &Config{
		Env:                           env,
		CreateInvestmentQueueURL:      os.Getenv("CREATE_INVESTMENT_QUEUE_URL"),
		CalculateAveragePriceQueueURL: os.Getenv("CALCULATE_AVERAGE_PRICE_QUEUE_URL"),
		Cloudflare: CloudflareConfig{
			AccountId:           os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
			InvestmentTrackDbId: os.Getenv("CLOUDFLARE_DB_ID"),
			ApiKey:              os.Getenv("CLOUDFLARE_API_KEY"),
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
