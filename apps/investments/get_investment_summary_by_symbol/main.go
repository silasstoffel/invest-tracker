package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	http_helper "github.com/silasstoffel/invest-tracker/apps/shared/http_helpers"
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	type CustomResponse struct {
		Code    string `json:"raw_code"`
		Message string `json:"raw_message"`
	}

	data := CustomResponse{
		Code:    "SUCCESS",
		Message: "Investment summary retrieved successfully",
	}

	return http_helper.JsonResponse(data), nil
}

func main() {
	lambda.Start(Handler)
}
