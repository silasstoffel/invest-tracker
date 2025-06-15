package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Response events.APIGatewayProxyResponse

var (
	sqsClient *sqs.Client
	queueURL  string
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("Failure to load aws config: %v", err))
	}
	sqsClient = sqs.NewFromConfig(cfg)
	queueURL = os.Getenv("CREATE_INVESTMENT_QUEUE_URL")
}

func Handler(ctx context.Context) (Response, error) {
	messageBody := map[string]interface{}{
		"action":  "create_investment",
		"payload": "dados fictícios para teste",
	}

	bodyJSON, err := json.Marshal(messageBody)
	if err != nil {
		return Response{StatusCode: 500}, fmt.Errorf("erro ao serializar payload: %w", err)
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(bodyJSON)),
	})
	if err != nil {
		return Response{StatusCode: 500}, fmt.Errorf("erro ao enviar mensagem SQS: %w", err)
	}

	// Resposta HTTP da Lambda
	respBody, _ := json.Marshal(map[string]string{
		"message": "Mensagem enviada para SQS com sucesso!",
	})

	var buf bytes.Buffer
	json.HTMLEscape(&buf, respBody)

	return Response{
		StatusCode:      202,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
