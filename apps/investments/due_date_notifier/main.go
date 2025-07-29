package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/silasstoffel/invest-tracker/apps/shared/telegram"
	appConfig "github.com/silasstoffel/invest-tracker/config"
)

func Handler() error {
	env := appConfig.NewConfigFromEnvVars()
	if env.Env == "production" {
		prefix := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
		tb := telegram.NewTelegramBot(env)
		tb.SendMessage(fmt.Sprintf("*[%s] Testing lambda cron*", prefix))
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
