package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	investment_summary_core "github.com/silasstoffel/invest-tracker/apps/investments_summary/core"
)

func handleMessage(msg string) error {
	var data investment_summary_core.InvestmentCreatedInput
	err := json.Unmarshal([]byte(msg), &data)

	if err != nil {
		log.Printf("Failure to convert JSON to struct")
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
