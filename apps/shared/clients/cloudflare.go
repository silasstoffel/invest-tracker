package client

import (
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

func (c *Client) InitCloudflare(apiKey string) {
	c.CloudflareClient = cloudflare.NewClient(option.WithAPIToken(apiKey))
}

func (c *Client) CreateCloudFlareClient(apiKey string) *cloudflare.Client {
	cfClient := cloudflare.NewClient(
		option.WithAPIToken(apiKey),
	)
	return cfClient
}
