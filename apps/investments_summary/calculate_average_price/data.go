package main

import (
	"github.com/cloudflare/cloudflare-go/v4"
	appConfig "github.com/silasstoffel/invest-tracker/config"
)

var (
	cfClient *cloudflare.Client
	env      *appConfig.Config
)
