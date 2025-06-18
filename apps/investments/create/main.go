package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func createInvestment(input string) error {
	// Simulate creating an investment
	log.Println("Creating investment...")
	// Here you would typically interact with a database or another service
	return nil

}

func Handler(ctx context.Context, sqsEvent events.SQSEvent) (map[string]interface{}, error) {
	batchItemFailures := []map[string]interface{}{}

	for _, message := range sqsEvent.Records {
		log.Printf("Received message: %s", message.Body)

		if message.Body == "" {
			log.Printf("Message with ID %s is empty, skipping", message.MessageId)
			batchItemFailures = append(batchItemFailures, map[string]interface{}{"itemIdentifier": message.MessageId})
		}

		err := createInvestment(message.Body)

		if err != nil {
			log.Printf("Error processing message %s: %v", message.MessageId, err)
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
