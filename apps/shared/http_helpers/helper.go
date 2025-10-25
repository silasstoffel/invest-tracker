package http_helper

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type JsonResponseOptions struct {
	Headers    map[string]string
	StatusCode int
}

func Response(body string, statusCode int, headers map[string]string) events.APIGatewayProxyResponse {
	status := 200
	if statusCode >= 200 && statusCode < 600 {
		status = statusCode
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      status,
		IsBase64Encoded: false,
		Body:            body,
		Headers:         headers,
	}
}

func JsonResponse(data interface{}, opts ...JsonResponseOptions) events.APIGatewayProxyResponse {
	statusCode := 200
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	if len(opts) > 0 {
		option := opts[0]

		if option.StatusCode != 0 {
			statusCode = option.StatusCode
		}

		if option.Headers != nil {
			for k, v := range option.Headers {
				headers[k] = v
			}
		}
	}
	encoded, err := json.Marshal(data)
	if err != nil {
		return Response(`{"code": "INTERNAL_ERROR", "message":"failed to marshal JSON"}`, 500, headers)
	}

	return Response(string(encoded), statusCode, headers)
}
