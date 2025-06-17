package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, sqsEvent events.SQSEvent) (map[string]interface{}, error) {
	batchItemFailures := []map[string]interface{}{}

	for _, message := range sqsEvent.Records {
		log.Printf("Received message: %s", message.Body)

		if message.Body == "" {
			log.Printf("Message with ID %s is empty, skipping", message.MessageId)
			batchItemFailures = append(batchItemFailures, map[string]interface{}{"itemIdentifier": message.MessageId})
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
