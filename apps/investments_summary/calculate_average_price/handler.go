package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/cloudflare/cloudflare-go/v4"
	investment_summary_core "github.com/silasstoffel/invest-tracker/apps/investments_summary/core"
	client "github.com/silasstoffel/invest-tracker/apps/shared/clients"
	appConfig "github.com/silasstoffel/invest-tracker/config"
)

var (
	cfClient *cloudflare.Client
	clients  *client.Client
	ctx      context.Context
)

func init() {
	env := appConfig.NewConfigFromEnvVars()
	ctx = context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		m := fmt.Sprintf("Failure to load aws config: %v", err)
		log.Println(m)
		panic(m)
	}
	ssmClient := ssm.NewFromConfig(cfg)
	ssmOutput, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &env.Cloudflare.ApiKeyParamName,
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		m := fmt.Sprintf("Failure to load data from parameter store: %v", err)
		log.Println(m)
		panic(m)
	}

	clients = client.CreateNewClients()
	clients.InitCloudflare(*ssmOutput.Parameter.Value)
	cfClient = clients.CloudflareClient
}

func GetInvestments(input investment_summary_core.InvestmentCreatedInput) {
	params := []string{
		input.Symbol,
		input.Type,
	}
	command := "select * from investments where symbol = ? and type = ?"
	res, err := clients.ExecuteD1RawQuery(ctx, cfClient, command, params)

	if err != nil {
		return err
	}

	res.Result.Rows

}
