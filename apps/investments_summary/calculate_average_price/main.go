package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	cryptoRand "crypto/rand"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/d1"
	"github.com/oklog/ulid/v2"
	investment_summary_core "github.com/silasstoffel/invest-tracker/apps/investments_summary/core"
	client "github.com/silasstoffel/invest-tracker/apps/shared/clients"
	appConfig "github.com/silasstoffel/invest-tracker/config"
)

type GetSummarizedInvestmentOutput struct {
	ID           string  `json:"id"`
	Quantity     float64 `json:"quantity"`
	AveragePrice float64 `json:"average_price"`
	TotalValue   float64 `json:"total_value"`
	Cost         float64 `json:"cost"`
}

type UpdateSummarizedInvestmentInput struct {
	ID           string
	Quantity     float64
	AveragePrice float64
	TotalValue   float64
	Cost         float64
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

func createId() string {
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)
	t := time.Now().UTC()

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func getSummarizedInvestment(input investment_summary_core.InvestmentCreatedInput) (GetSummarizedInvestmentOutput, error) {
	var params []string
	command := ""

	if input.OperationType == "sell" && input.Type == "bond" {
		command = "select id, quantity, average_price, total_value, cost from investments_summary where investment_id = ? limit 1"
		params = []string{
			input.SellInvestmentId,
		}
	} else {
		params = []string{
			input.Symbol,
			input.Type,
			input.Brokerage,
		}
		command = "select id, quantity, average_price, total_value, cost from investments_summary where symbol = ? and type = ? and brokerage = ? limit 1"
	}

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
		jsonInput, err := json.Marshal(rawRow)
		if err != nil {
			return GetSummarizedInvestmentOutput{}, fmt.Errorf("failure to convert interface to json: %w", err)
		}

		var output GetSummarizedInvestmentOutput
		if err := json.Unmarshal(jsonInput, &output); err != nil {
			return GetSummarizedInvestmentOutput{}, fmt.Errorf("failure to convert json to struct: %w", err)
		}

		return output, nil
	}

	return GetSummarizedInvestmentOutput{}, errors.New("summarized investment not found")
}

func createSummarizedInvestment(input investment_summary_core.InvestmentCreatedInput) (string, error) {
	if input.OperationType != "buy" && input.OperationType != "sell" {
		return "", errors.New("operation type must be 'buy' or 'sell'")
	}

	command := `insert into investments_summary(
		id, investment_id, brokerage, type, symbol,
		quantity, average_price, total_value, cost,
		redemption_policy_type, created_at, updated_at, last_operation_date {add_column_name}
	) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? {add_column_value})`

	avgPrice := input.TotalValue / float64(input.Quantity)
	createdAt := time.Now().Format(time.RFC3339)
	id := createId()

	params := []string{
		id,
		input.ID,
		input.Brokerage,
		input.Type,
		input.Symbol,
		fmt.Sprintf("%v", input.Quantity),
		fmt.Sprintf("%v", avgPrice),
		fmt.Sprintf("%v", input.TotalValue),
		fmt.Sprintf("%v", input.Cost),
		input.RedemptionPolicyType,
		createdAt,
		createdAt,
		input.OperationDate,
	}

	if input.BondIndex != "" {
		command = strings.Replace(command, "{add_column_name}", ", bond_index, bond_rate", 1)
		command = strings.Replace(command, "{add_column_value}", ",?,?", 1)
		params = append(params, input.BondIndex, fmt.Sprintf("%v", input.BondRate))
	} else {
		command = strings.Replace(command, "{add_column_name}", "", 1)
		command = strings.Replace(command, "{add_column_value}", "", 1)
	}

	_, err := cfClient.D1.Database.Raw(
		context.TODO(),
		env.Cloudflare.InvestmentTrackDbId,
		d1.DatabaseRawParams{
			AccountID: cloudflare.F(env.Cloudflare.AccountId),
			Sql:       cloudflare.F(command),
			Params:    cloudflare.F(params),
		},
	)

	if err != nil {
		m := fmt.Sprintf("Error executing command on cloudflare. Detail: %v", err)
		log.Printf(m)
		log.Printf("Command: %s", command)
		log.Printf("Params: %v", params)
		return "", errors.New(m)
	}

	saveHistory(id)
	return id, nil
}

func updateSummarizedInvestment(currentPosition UpdateSummarizedInvestmentInput, createdInvestment investment_summary_core.InvestmentCreatedInput) error {
	if createdInvestment.OperationType != "buy" && createdInvestment.OperationType != "sell" {
		return errors.New("operation type must be 'buy' or 'sell'")
	}

	quantity := currentPosition.Quantity
	averagePrice := currentPosition.AveragePrice
	totalValue := currentPosition.TotalValue
	costs := currentPosition.Cost

	if createdInvestment.OperationType == "sell" {
		if quantity == createdInvestment.Quantity {
			quantity = 0
			averagePrice = 0
			totalValue = 0
			costs = 0
		} else if createdInvestment.Type == "bond" {
			totalValue -= createdInvestment.TotalValue
			averagePrice = totalValue / quantity
		} else {
			// average price does not change when the operation type is sell
			quantity -= createdInvestment.Quantity
			// total value is reduced by the average price times the quantity sold
			totalValue -= (averagePrice * float64(createdInvestment.Quantity))
			costs -= createdInvestment.Cost
		}
	} else {
		quantity += createdInvestment.Quantity
		totalValue += createdInvestment.TotalValue
		averagePrice = totalValue / float64(quantity)
		costs += createdInvestment.Cost
	}

	params := []string{
		fmt.Sprintf("%f", quantity),
		fmt.Sprintf("%f", averagePrice),
		fmt.Sprintf("%f", totalValue),
		fmt.Sprintf("%f", costs),
		time.Now().UTC().Format(time.RFC3339),
		createdInvestment.OperationDate,
		currentPosition.ID,
	}

	command := "update investments_summary set quantity = ?, average_price = ?, total_value = ?, cost = ?, updated_at = ?, last_operation_date = ? where id = ?"

	_, err := cfClient.D1.Database.Query(context.TODO(), env.Cloudflare.InvestmentTrackDbId, d1.DatabaseQueryParams{
		AccountID: cloudflare.F(env.Cloudflare.AccountId),
		Sql:       cloudflare.F(command),
		Params:    cloudflare.F(params),
	})

	if err != nil {
		return fmt.Errorf("failure to update summarized investment: %w", err)
	}

	saveHistory(currentPosition.ID)

	return nil
}

func saveHistory(summarizedInvestmentId string) error {
	command := `INSERT INTO investments_summary_history(
    investment_id,
    last_operation_date,
    operation_month,
    operation_year, 
    brokerage,
    type,
    symbol,
    bond_index,
    bond_rate,
    quantity,
    average_price,
    total_value,
    market_value,
    cost,
    redemption_policy_type,
    due_date,
    investment_summary_id   
) SELECT 
    investment_id,
    last_operation_date,
    strftime('%m', last_operation_date) as operation_month,
    strftime('%Y', last_operation_date) as operation_year,
    brokerage,
    type,
    symbol,
    bond_index,
    bond_rate,
    quantity,
    average_price,
    total_value,
    market_value,
    cost,
    redemption_policy_type,
    due_date,
    id as investment_summary_id 
  FROM investments_summary
  WHERE id = ?`

	params := []string{
		summarizedInvestmentId,
	}
	_, err := cfClient.D1.Database.Raw(
		context.TODO(),
		env.Cloudflare.InvestmentTrackDbId,
		d1.DatabaseRawParams{
			AccountID: cloudflare.F(env.Cloudflare.AccountId),
			Sql:       cloudflare.F(command),
			Params:    cloudflare.F(params),
		},
	)

	if err != nil {
		m := fmt.Sprintf("[save-history]: error when save history. Detail: %v", err)
		log.Printf(m)
		log.Printf("Command: %s", command)
		log.Printf("Params: %v", params)
		return errors.New(m)
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
			summarizedId, err := createSummarizedInvestment(input)

			if err != nil {
				log.Printf("Failure to create summarized investment: %v", err)
				return err
			}
			log.Printf("Created summarized investment with ID: %s - Symbol: %s", summarizedId, input.Symbol)

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
		Cost:         summarized.Cost,
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
