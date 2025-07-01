package client

import (
	"github.com/cloudflare/cloudflare-go/v4"
)

type Client struct {
	CloudflareClient *cloudflare.Client
}

func CreateNewClient() *Client {
	return &Client{}
}
