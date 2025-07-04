package main

import (
	"context"
	"log"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/d1"
	client "github.com/silasstoffel/invest-tracker/apps/shared/clients"
)

type item struct {
	Symbol string `json:"symbol"`
}

func main() {
	clients := client.CreateNewClients()
	clients.InitCloudflare("EzzZwypF4AO5Dm94WAvXr58EzL8baMXhcmveJ4sr")

	cfClient := clients.CloudflareClient

	dbId := "e256f11a-1b2e-4add-bcec-36697eb80eec"
	params := []string{}

	res, err := cfClient.D1.Database.Query(context.TODO(), dbId, d1.DatabaseQueryParams{
		AccountID: cloudflare.F("b3262b1ed9d85abab8621e13d0aba2aa"),
		Sql:       cloudflare.F("select * from investments where 1=1 limit 1"),
		Params:    cloudflare.F(params),
	})

	if err != nil {
		log.Fatalf("query failed: %v", err)
	}

	for _, queryResult := range res.Result {
		for i, rawRow := range queryResult.Results {
			row, ok := rawRow.(map[string]interface{})
			if !ok {
				log.Printf("Linha %d: formato inesperado: %T", i+1, rawRow)
				continue
			}

			typeInvesment, _ := row["type"].(string)
			symbol, _ := row["symbol"].(string)
			totalValue, _ := row["total_value"].(float64)

			log.Print(typeInvesment, symbol, totalValue)
			/*
				log.Printf("Linha %d:", i+1)
				for col, val := range row {
					log.Printf("  %s: %v", col, val)
				}
			*/
		}
	}

}
