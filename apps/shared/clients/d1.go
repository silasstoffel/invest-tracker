package client

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/d1"
	"github.com/cloudflare/cloudflare-go/v4/packages/pagination"
)

func (c *Client) ExecuteD1RawQuery(ctx context.Context, cloudflareClient *cloudflare.Client, command string, params []string) (res *pagination.SinglePage[d1.DatabaseRawResponse], err error) {
	res, err = cloudflareClient.D1.Database.Raw(
		ctx,
		c.config.Cloudflare.InvestmentTrackDbId,
		d1.DatabaseRawParams{
			AccountID: cloudflare.F(c.config.Cloudflare.AccountId),
			Sql:       cloudflare.F(command),
			Params:    cloudflare.F(params),
		},
	)

	return res, err
}
