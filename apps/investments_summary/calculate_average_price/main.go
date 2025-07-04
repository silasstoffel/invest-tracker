package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/d1"
	investment_summary_core "github.com/silasstoffel/invest-tracker/apps/investments_summary/core"
	client "github.com/silasstoffel/invest-tracker/apps/shared/clients"
	appConfig "github.com/silasstoffel/invest-tracker/config"
)

type GetSummarizedInvestmentOutput struct {
	ID           string
	Quantity     int
	AveragePrice float64
	TotalValue   float64
	Costs        float64
}

type UpdateSummarizedInvestmentInput struct {
	ID           string
	Quantity     int
	AveragePrice float64
	TotalValue   float64
	Costs        float64
}

var (
	cfClient *cloudflare.Client
	clients  *client.Client
	ctx      context.Context
	env      *appConfig.Config
)

func init() {
	env = appConfig.NewConfigFromEnvVars()
	ctx = context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		m := fmt.Sprintf("Failure to load aws config: %v", err)
		log.Println(m)
		panic(m)
	}
	ssmClient := ssm.NewFromConfig(cfg)
	ssmOutput, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &env.Cloudflare.ApiKeyParamName,
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		m := fmt.Sprintf("Failure to load data from parameter store: %v", err)
		log.Println(m)
		panic(m)
	}

	clients = client.CreateNewClients()
	clients.InitCloudflare(*ssmOutput.Parameter.Value)
	cfClient = clients.CloudflareClient
}

func getSummarizedInvestment(input investment_summary_core.InvestmentCreatedInput) (GetSummarizedInvestmentOutput, error) {
	params := []string{
		input.Symbol,
		input.Type,
		input.Brokerage,
	}
	command := "select id, quantity, average_price, total_value, costs from investments_summary where symbol = ? and type = ? and brokerage = ? limit 1"

	res, err := cfClient.D1.Database.Query(context.TODO(), env.Cloudflare.InvestmentTrackDbId, d1.DatabaseQueryParams{
		AccountID: cloudflare.F(env.Cloudflare.AccountId),
		Sql:       cloudflare.F(command),
		Params:    cloudflare.F(params),
	})

	if err != nil {
		return GetSummarizedInvestmentOutput{}, err
	}

	if len(res.Result) > 0 && len(res.Result[0].Results) > 0 {
		rawRow := res.Result[0].Results[0]
		row, ok := rawRow.(map[string]interface{})
		if !ok {
			return GetSummarizedInvestmentOutput{}, errors.New("failure when casting cloudflare query result")
		}

		id, _ := row["id"].(string)
		quantity, _ := row["quantity"].(int)
		averagePrice, _ := row["average_price"].(float64)
		totalValue, _ := row["total_value"].(float64)
		costs, _ := row["total_value"].(float64)

		return GetSummarizedInvestmentOutput{
			ID:           id,
			Quantity:     quantity,
			AveragePrice: averagePrice,
			TotalValue:   totalValue,
			Costs:        costs,
		}, nil
	}

	return GetSummarizedInvestmentOutput{}, errors.New("summarized investment not found")
}

func createSummarizedInvestment(input investment_summary_core.InvestmentCreatedInput) (GetSummarizedInvestmentOutput, error) {
	if input.OperationType != "buy" && input.OperationType != "sell" {
		return GetSummarizedInvestmentOutput{}, errors.New("operation type must be 'buy' or 'sell'")
	}

	return GetSummarizedInvestmentOutput{}, nil
}

func updateSummarizedInvestment(input UpdateSummarizedInvestmentInput, createdInvestment investment_summary_core.InvestmentCreatedInput) error {
	if createdInvestment.OperationType != "buy" && createdInvestment.OperationType != "sell" {
		return errors.New("operation type must be 'buy' or 'sell'")
	}

	quantity := input.Quantity
	averagePrice := input.AveragePrice
	totalValue := input.TotalValue
	costs := input.Costs

	if createdInvestment.OperationType == "sell" {
		quantity -= createdInvestment.Quantity
		totalValue -= createdInvestment.TotalValue
		costs -= createdInvestment.Cost
	} else {
		quantity += createdInvestment.Quantity
		totalValue += createdInvestment.TotalValue
		averagePrice = totalValue / float64(quantity)
		costs += createdInvestment.Cost
	}

	params := []string{
		fmt.Sprintf("%d", quantity),
		fmt.Sprintf("%f", averagePrice),
		fmt.Sprintf("%f", totalValue),
		fmt.Sprintf("%f", costs),
		time.Now().UTC().Format(time.RFC3339),
		input.ID,
	}

	command := "update investments_summary set quantity = ?, average_price = ?, total_value = ?, costs = ?, updated_at = ? where id = ?"

	_, err := cfClient.D1.Database.Query(context.TODO(), env.Cloudflare.InvestmentTrackDbId, d1.DatabaseQueryParams{
		AccountID: cloudflare.F(env.Cloudflare.AccountId),
		Sql:       cloudflare.F(command),
		Params:    cloudflare.F(params),
	})

	if err != nil {
		return fmt.Errorf("failure to update summarized investment: %w", err)
	}

	return nil
}

func handleMessage(msg string) error {
	var input investment_summary_core.InvestmentCreatedInput
	err := json.Unmarshal([]byte(msg), &input)

	if err != nil {
		log.Printf("Failure to convert JSON to struct")
		return err
	}

	summarized, err := getSummarizedInvestment(input)

	if err != nil {
		if err.Error() == "summarized investment not found" {
			log.Printf("Summarized investment not found for symbol %s. It does need to create it", input.Symbol)
			created, err := createSummarizedInvestment(input)

			if err != nil {
				log.Printf("Failure to create summarized investment: %v", err)
				return err
			}
			log.Printf("Created summarized investment with ID: %s - Symbol: %s", created.ID, input.Symbol)

			return nil
		}

		log.Printf("Failure to get summarized investment: %v", err)
		return err
	}

	err = updateSummarizedInvestment(UpdateSummarizedInvestmentInput{
		ID:           summarized.ID,
		Quantity:     summarized.Quantity,
		AveragePrice: summarized.AveragePrice,
		TotalValue:   summarized.TotalValue,
		Costs:        summarized.Costs,
	}, input)

	if err != nil {
		log.Printf("Failure to update summarized investment: %v", err)
		return err
	}

	return nil
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	batchItemFailures := []events.SQSBatchItemFailure{}

	for _, message := range sqsEvent.Records {
		err := handleMessage(message.Body)

		if err != nil {
			log.Printf("Error processing message %s: %v", message.MessageId, err)
			log.Printf("Received message: %s", message.Body)
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: message.MessageId})
			continue
		}
	}

	return events.SQSEventResponse{
		BatchItemFailures: batchItemFailures,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
