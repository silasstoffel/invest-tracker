package main

import (
	"context"
	cryptoRand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/d1"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/oklog/ulid/v2"
	investment_core "github.com/silasstoffel/invest-tracker/apps/investments/core"
	appConfig "github.com/silasstoffel/invest-tracker/config"
)

var (
	cfClient  *cloudflare.Client
	env       *appConfig.Config
	sqsClient *sqs.Client
	queueURL  string
)

func init() {
	env = appConfig.NewConfigFromEnvVars()
	ctx := context.Background()
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

	cfClient = cloudflare.NewClient(
		option.WithAPIToken(*ssmOutput.Parameter.Value),
	)

	sqsClient = sqs.NewFromConfig(cfg)
	config := appConfig.NewConfigFromEnvVars()
	queueURL = config.CalculateAveragePriceQueueURL
}

func createId() string {
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)
	t := time.Now().UTC()

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func saveInvestment(entity investment_core.InvestmentEntity) error {
	command := `INSERT INTO investments (
		id, type, symbol, quantity, unit_price, total_value, cost, operation_type, operation_date,
		operation_year, operation_month, due_date, created_at, updated_at, brokerage, note, redemption_policy_type, sell_investment_id {add_column_name}) VALUES (
			?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?{add_column_value})`

	params := []string{
		entity.ID,
		entity.Type,
		entity.Symbol,
		fmt.Sprintf("%v", entity.Quantity),
		fmt.Sprintf("%v", entity.UnitPrice),
		fmt.Sprintf("%v", entity.TotalValue),
		fmt.Sprintf("%v", entity.Cost),
		entity.OperationType,
		entity.OperationDate,
		fmt.Sprintf("%v", entity.OperationYear),
		fmt.Sprintf("%v", entity.OperationMonth),
		entity.DueDate,
		entity.CreatedAt.Format(time.RFC3339),
		entity.UpdatedAt.Format(time.RFC3339),
		entity.Brokerage,
		entity.Note,
		entity.RedemptionPolicyType,
		entity.SellInvestmentId,
	}

	if entity.BondIndex != "" {
		command = strings.Replace(command, "{add_column_name}", ", bond_index, bond_rate", 1)
		command = strings.Replace(command, "{add_column_value}", ",?,?", 1)
		params = append(params, entity.BondIndex, fmt.Sprintf("%v", entity.BondRate))
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
		return errors.New(m)
	}

	return nil
}

func createInvestment(input string) (investment_core.InvestmentEntity, error) {
	var data investment_core.CreateInvestmentInput
	err := json.Unmarshal([]byte(input), &data)

	if err != nil {
		log.Println("Failure to convert json input to create investment input. Detail: ", input)
		return investment_core.InvestmentEntity{}, err
	}

	od, _ := time.Parse("2006-01-02", data.OperationDate)
	entity := investment_core.InvestmentEntity{
		ID:                   createId(),
		Type:                 data.Type,
		Symbol:               data.Symbol,
		BondIndex:            data.BondIndex,
		BondRate:             data.BondRate,
		Quantity:             data.Quantity,
		UnitPrice:            data.TotalValue / data.Quantity,
		TotalValue:           data.TotalValue,
		Cost:                 data.Cost,
		OperationType:        data.OperationType,
		OperationDate:        data.OperationDate,
		OperationYear:        od.Year(),
		OperationMonth:       int(od.Month()),
		DueDate:              data.DueDate,
		Brokerage:            data.Brokerage,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
		RedemptionPolicyType: data.RedemptionPolicyType,
		Note:                 data.Note,
		SellInvestmentId:     data.SellInvestmentId,
	}

	err = saveInvestment(entity)
	if err != nil {
		log.Printf("Error saving investment: %v", err)
		return investment_core.InvestmentEntity{}, fmt.Errorf("error saving investment: %w", err)
	}

	log.Println("Investment created successfully. ID:", entity.ID, " Symbol:", entity.Symbol, " Type:", entity.Type, " Total Value:", entity.TotalValue)
	return entity, nil
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	batchItemFailures := []events.SQSBatchItemFailure{}

	for _, message := range sqsEvent.Records {

		if message.Body == "" {
			log.Printf("Message with ID %s is empty, skipping", message.MessageId)
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: message.MessageId})
			continue
		}

		entity, err := createInvestment(message.Body)

		if err != nil {
			log.Printf("Error processing message %s: %v", message.MessageId, err)
			log.Printf("Received message: %s", message.Body)
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{ItemIdentifier: message.MessageId})
			continue
		}

		messageContent, err := json.Marshal(entity)
		if err != nil {
			m := "Failure when converting entity message to json. Message was not sent to recalculate average price. Detail: %v"
			log.Printf(m, err)
			continue
		}
		log.Printf("Sending message to calculate-average-price-env-queue.fifo: %s", string(messageContent))
		_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:               aws.String(queueURL),
			MessageBody:            aws.String(string(messageContent)),
			MessageDeduplicationId: aws.String(entity.ID),
			MessageGroupId:         aws.String(strings.ReplaceAll(entity.Symbol, " ", "_")),
		})

		if err != nil {
			m := "Failure to send message to calculate-average-price-env-queue.fifo. Detail: %v"
			log.Printf(m, err)
		}
	}

	return events.SQSEventResponse{
		BatchItemFailures: batchItemFailures,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
