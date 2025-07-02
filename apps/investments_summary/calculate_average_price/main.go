package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleMessage() error {
	return nil
}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	batchItemFailures := []events.SQSBatchItemFailure{}

	for _, message := range sqsEvent.Records {
		err := handleMessage()

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
