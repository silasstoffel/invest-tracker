# Invest Track

## About

Serverless app to manage my personal investments.

![Diagram](./docs/invest-track.png)

## Requirements

 - [Go 1.24.3](https://go.dev/)
 - [Serverless framework](https://www.serverless.com/)
 - [Cloudflare D1](https://developers.cloudflare.com/d1/)

## Deploy on AWS

Make sure you have done AWS CLI setup.

```shell
# development account
make deploy-dev

# production account
make deploy-prod
```