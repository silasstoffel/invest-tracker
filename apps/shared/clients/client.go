package client

import (
	"github.com/cloudflare/cloudflare-go/v4"
	appConfig "github.com/silasstoffel/invest-tracker/config"
)

type Client struct {
	CloudflareClient *cloudflare.Client
	config           *appConfig.Config
}

func CreateNewClients() *Client {
	return &Client{
		config: appConfig.NewConfigFromEnvVars(),
	}
}
