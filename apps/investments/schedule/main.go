package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	investment_core "github.com/silasstoffel/invest-tracker/apps/investments/core"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Response events.APIGatewayProxyResponse

var (
	sqsClient       *sqs.Client
	queueURL        string
	responseHeaders = map[string]string{
		"Content-Type": "application/json",
	}
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("Failure to load aws config: %v", err))
	}
	sqsClient = sqs.NewFromConfig(cfg)
	queueURL = os.Getenv("CREATE_INVESTMENT_QUEUE_URL")
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var input investment_core.CreateInvestmentInput
	err := json.Unmarshal([]byte(request.Body), &input)

	if err != nil {
		message := "Unsupported input format"
		respBody, _ := json.Marshal(map[string]string{
			"message": message,
			"code":    "INVALID_INPUT",
		})

		log.Println(string(respBody))

		return Response{
			StatusCode: 400,
			Headers:    responseHeaders,
			Body:       string(respBody),
		}, fmt.Errorf("%s: %w", message, err)
	}

	messageContent, err := json.Marshal(input)
	if err != nil {
		message := "Failure to convert message input to JSON"
		respBody, _ := json.Marshal(map[string]string{
			"message": message,
			"code":    "INTERNAL_ERROR",
		})

		log.Printf("%s: %v", message, err)

		return Response{
			StatusCode: 400,
			Headers:    responseHeaders,
			Body:       string(respBody),
		}, fmt.Errorf("%s: %w", message, err)
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(messageContent)),
	})

	if err != nil {
		const message = "Failed to send message to SQS"
		respBody, _ := json.Marshal(map[string]string{
			"message": message,
			"code":    "INTEGRATION_ERROR",
		})
		return Response{
			StatusCode: 400,
			Headers:    responseHeaders,
			Body:       string(respBody),
		}, fmt.Errorf("%s: %w", message, err)
	}

	response, _ := json.Marshal(map[string]string{
		"message": "Investment will be created as a soon as possible!",
	})

	return Response{
		StatusCode:      202,
		IsBase64Encoded: false,
		Body:            string(response),
		Headers:         responseHeaders,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
