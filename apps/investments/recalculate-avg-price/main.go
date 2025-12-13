package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	appConfig "github.com/silasstoffel/invest-tracker/config"
)

var (
	sqsClient *sqs.Client
	queueURL  string
)

/* Schema of each item in the JSON array:
{
  "id": "01JYAMF0KA2CGV5MPD0BNHCRH4",
  "type": "bond",
  "symbol": "LCA PRE BTG",
  "bondIndex": "prefixed",
  "bondRate": 11,
  "quantity": 1,
  "unitPrice": 942,
  "totalValue": 942,
  "cost": 0,
  "operationType": "buy",
  "operationDate": "2022-02-21",
  "operationYear": 2022,
  "operationMonth": 2,
  "dueDate": "2028-01-19",
  "brokerage": "banco inter",
  "note": "inter_import",
  "redemptionPolicyType": "at_maturity",
  "sellInvestmentId": null,
  "createdAt": "2025-06-22T01:36:21Z",
  "updatedAt": "2025-06-22T01:36:21Z"
}
*/

var strJson = `[]`

func init() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		m := fmt.Sprintf("Failure to load aws config: %v", err)
		log.Println(m)
		panic(m)
	}

	sqsClient = sqs.NewFromConfig(cfg)
	config := appConfig.NewConfigFromEnvVars()
	queueURL = config.CalculateAveragePriceQueueURL
}

func Handler(ctx context.Context) error {
	items := []map[string]interface{}{}
	err := json.Unmarshal([]byte(strJson), &items)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	for _, item := range items {
		symbol := item["symbol"].(string)
		id := item["id"].(string)

		itemJSON, err := json.Marshal(item)
		if err != nil {
			log.Printf("Error marshaling item: %v", err)
			continue
		}

		_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:               aws.String(queueURL),
			MessageBody:            aws.String(string(itemJSON)),
			MessageDeduplicationId: aws.String(id),
			MessageGroupId:         aws.String(strings.ReplaceAll(symbol, " ", "_")),
		})

		if err != nil {
			m := "Failure to send message to calculate-average-price-env-queue.fifo. Detail: %v"
			log.Printf(m, err)
		}
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
