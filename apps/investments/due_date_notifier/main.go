package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/d1"
	client "github.com/silasstoffel/invest-tracker/apps/shared/clients"
	"github.com/silasstoffel/invest-tracker/apps/shared/telegram"
	appConfig "github.com/silasstoffel/invest-tracker/config"
)

var (
	cfClient *cloudflare.Client
	clients  *client.Client
	env      *appConfig.Config
)

func init() {
	env = appConfig.NewConfigFromEnvVars()

	clients = client.CreateNewClients()
	clients.InitCloudflare(env.Cloudflare.ApiKey)
	cfClient = clients.CloudflareClient
}

type GetInvestmentOutput struct {
	Brokerage  string  `json:"brokerage"`
	Type       string  `json:"type"`
	Symbol     string  `json:"symbol"`
	DueDate    string  `json:"due_date"`
	TotalValue float64 `json:"total_value"`
}

func getInvestments() ([]GetInvestmentOutput, error) {
	today := time.Now()
	f := "2006-01-02"
	startAt := today.Format(f)
	finishAt := today.AddDate(0, 0, 7).Format(f)

	params := []string{startAt, finishAt}
	command := `SELECT ins.brokerage, ins.type, ins.symbol, ins.due_date, ins.total_value 
				FROM investments_summary ins 
				WHERE ins.due_date <> ''
				  AND ins.due_date BETWEEN ? AND ?
				ORDER BY ins.due_date`

	res, err := cfClient.D1.Database.Query(context.TODO(), env.Cloudflare.InvestmentTrackDbId, d1.DatabaseQueryParams{
		AccountID: cloudflare.F(env.Cloudflare.AccountId),
		Sql:       cloudflare.F(command),
		Params:    cloudflare.F(params),
	})

	if err != nil {
		log.Printf("Failure do execute command: %s. Params: %v Detail: %v", command, params, err)
		return nil, err
	}

	outputs := []GetInvestmentOutput{}

	if len(res.Result) == 0 || len(res.Result[0].Results) == 0 {
		log.Printf("getInvestments: no results found for command")
		return outputs, nil
	}

	for _, row := range res.Result[0].Results {
		jsonRow, err := json.Marshal(row)
		if err != nil {
			log.Printf("Failed to marshal row to JSON. %v", err)
			return nil, fmt.Errorf("failed to marshal row to JSON: %w", err)
		}

		var output GetInvestmentOutput
		if err := json.Unmarshal(jsonRow, &output); err != nil {
			log.Printf("Failed to unmarshal row to struct. %v", err)
			return nil, fmt.Errorf("failed to unmarshal row to struct: %w", err)
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

func Handler() error {
	env := appConfig.NewConfigFromEnvVars()

	if env.Env != "production" {
		return nil
	}

	investments, err := getInvestments()
	prefix := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	tb := telegram.NewTelegramBot(env)

	if err != nil {
		tb.SendMessage(fmt.Sprintf("*[%s] Failure to read investments* ```%s```", prefix, err.Error()))
		return nil
	}

	counter := len(investments)
	if counter > 0 {
		jsonContent, _ := json.Marshal(investments)
		tb.SendMessage(fmt.Sprintf("*[%s] %d Investment(s) due this week* ```json %s```", prefix, counter, jsonContent))
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
