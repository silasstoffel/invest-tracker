package main

import (
	"context"
	cryptoRand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/d1"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/oklog/ulid/v2"
	investment_core "github.com/silasstoffel/invest-tracker/apps/investments/core"
)

var cfClient *cloudflare.Client

func init() {
	cfClient = cloudflare.NewClient(
		option.WithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN_DEV")),
	)
}

func createId() string {
	entropy := ulid.Monotonic(cryptoRand.Reader, 0)
	t := time.Now().UTC()

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func saveInvestment(entity investment_core.InvestmentEntity) error {
	command := `INSERT INTO investments (
		id, type, symbol, quantity, unit_price, total_value, cost, operation_type, operation_date,
		operation_year, operation_month, due_date, created_at, updated_at{add_column_name}) VALUES (
			?,?,?,?,?,?,?,?,?,?,?,?,?,?{add_column_value})`

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
		os.Getenv("CLOUDFLARE_DB_ID_DEV"),
		d1.DatabaseRawParams{
			AccountID: cloudflare.F(os.Getenv("CLOUDFLARE_ACCOUNT_ID_DEV")),
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

func createInvestment(input string) error {
	var data investment_core.CreateInvestmentInput
	err := json.Unmarshal([]byte(input), &data)

	if err != nil {
		log.Println("Failure to convert json input to create investment input. Detail: ", input)
	}

	od, _ := time.Parse("2006-01-02", data.OperationDate)
	entity := investment_core.InvestmentEntity{
		ID:             createId(),
		Type:           data.Type,
		Symbol:         data.Symbol,
		BondIndex:      data.BondIndex,
		BondRate:       data.BondRate,
		Quantity:       data.Quantity,
		UnitPrice:      data.TotalValue / float64(data.Quantity),
		TotalValue:     data.TotalValue,
		Cost:           data.Cost,
		OperationType:  data.OperationType,
		OperationDate:  data.OperationDate,
		OperationYear:  od.Year(),
		OperationMonth: int(od.Month()),
		DueDate:        data.DueDate,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	err = saveInvestment(entity)
	if err != nil {
		log.Printf("Error saving investment: %v", err)
		return fmt.Errorf("error saving investment: %w", err)
	}

	log.Println("Investment created successfully. ID:", entity.ID, " Symbol:", entity.Symbol, " Type:", entity.Type, " Total Value:", entity.TotalValue)
	return nil

}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) (map[string]interface{}, error) {
	batchItemFailures := []map[string]interface{}{}

	for _, message := range sqsEvent.Records {

		if message.Body == "" {
			log.Printf("Message with ID %s is empty, skipping", message.MessageId)
			batchItemFailures = append(batchItemFailures, map[string]interface{}{"itemIdentifier": message.MessageId})
			continue
		}

		err := createInvestment(message.Body)

		if err != nil {
			log.Printf("Error processing message %s: %v", message.MessageId, err)
			log.Printf("Received message: %s", message.Body)
			batchItemFailures = append(batchItemFailures, map[string]interface{}{"itemIdentifier": message.MessageId})
			continue
		}
	}

	sqsBatchResponse := map[string]interface{}{
		"batchItemFailures": batchItemFailures,
	}

	return sqsBatchResponse, nil
}

func main() {
	lambda.Start(Handler)
}
