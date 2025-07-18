package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	investment_core "github.com/silasstoffel/invest-tracker/apps/investments/core"
	appConfig "github.com/silasstoffel/invest-tracker/config"

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
	config := appConfig.NewConfigFromEnvVars()
	queueURL = config.CreateInvestmentQueueURL
}

func checkInvestmentType(t string) error {
	switch t {
	case investment_core.FiiInvestmentType, investment_core.StockInvestmentType, investment_core.ReitInvestmentType, investment_core.BondInvestmentType, investment_core.EtfInvestmentType:
		return nil
	case "":
		return errors.New("investment type is required")
	default:
		return errors.New("invalid investment type")
	}
}

func checkRedemptionPolicyType(t string) error {
	switch t {
	case investment_core.AnyTimeRedemption, investment_core.AtMaturityRedemption, investment_core.HybridRedemption:
		return nil
	default:
		return errors.New("invalid redemption policy type")
	}
}

func checkOperationType(t string) error {
	switch t {
	case investment_core.BuyOperationType, investment_core.SellOperationType:
		return nil
	case "":
		return errors.New("operation type is required")
	default:
		return errors.New("invalid operation type")
	}
}

func validateInput(input investment_core.CreateInvestmentInput) error {
	investmentTypeError := checkInvestmentType(input.Type)
	if investmentTypeError != nil {
		return investmentTypeError
	}
	if input.Symbol == "" {
		return fmt.Errorf("investment symbol is required")
	}
	if input.Quantity <= 0 {
		return fmt.Errorf("investment quantity must be greater than zero")
	}
	if input.TotalValue <= 0 {
		return fmt.Errorf("investment total value must be greater than zero")
	}
	if input.Cost < 0 {
		return fmt.Errorf("investment cost cannot be negative")
	}
	checkOperationTypeErr := checkOperationType(input.OperationType)
	if checkOperationTypeErr != nil {
		return checkOperationTypeErr
	}
	if input.OperationDate == "" {
		return fmt.Errorf("operation date is required")
	} else {
		// Validate operation date format (YYYY-MM-DD)
		date := strings.Trim(input.OperationDate, "")
		if len(date) != 10 {
			return fmt.Errorf("operation date must be in the format YYYY-MM-DD")
		}
	}

	if input.DueDate != "" {
		// Validate due date format (YYYY-MM-DD)
		date := strings.Trim(input.DueDate, "")
		if len(date) != 10 {
			return fmt.Errorf("due date must be in the format YYYY-MM-DD")
		}
	}

	if input.Type == investment_core.BondInvestmentType {
		if input.BondIndex == "" {
			return fmt.Errorf("bond index is required for bond investments")
		}
		if (input.BondIndex == investment_core.BondIndexIPCA || input.BondIndex == investment_core.BondIndexPrefix) && input.BondRate < 0 {
			return fmt.Errorf("bond rate must be greater than zero")
		}
	}

	if input.RedemptionPolicyType != "" {
		if failure := checkRedemptionPolicyType(input.RedemptionPolicyType); failure != nil {
			return failure
		}
	}

	if input.OperationType == investment_core.SellOperationType && input.Type == investment_core.BondInvestmentType {
		sellInvestmentId := strings.Trim(input.SellInvestmentId, "")

		if sellInvestmentId == "" {
			return fmt.Errorf("sell investment ID is required for sell operation")
		}

		if len(sellInvestmentId) != 26 {
			return fmt.Errorf("sell investment ID must be 26 characters long")
		}
	}

	return nil
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
		}, nil
	}

	validateInputErr := validateInput(input)
	if validateInputErr != nil {
		message := validateInputErr.Error()
		respBody, _ := json.Marshal(map[string]string{
			"message": message,
			"code":    "INVALID_REQUEST",
		})
		log.Printf("%s: %v", message, validateInputErr)
		return Response{
			StatusCode: 400,
			Headers:    responseHeaders,
			Body:       string(respBody),
		}, nil
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
		}, nil
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
		}, nil
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
