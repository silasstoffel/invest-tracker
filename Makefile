.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	
	env GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -ldflags="-s -w" -o bin/bootstrap apps/investments/schedule/main.go
	cd ./bin && zip schedule-investment.zip bootstrap

	env GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -ldflags="-s -w" -o bin/bootstrap apps/investments/create/main.go
	cd ./bin && zip create-investment.zip bootstrap

	env GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -ldflags="-s -w" -o bin/bootstrap apps/investments/due_date_notifier/main.go
	cd ./bin && zip due-date-notifier.zip bootstrap

	env GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -ldflags="-s -w" -o bin/bootstrap apps/investments_summary/calculate_average_price/main.go
	cd ./bin && zip calculate-average-price.zip bootstrap

	env GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -ldflags="-s -w" -o bin/bootstrap apps/investments/get_investment_summary_by_symbol/main.go
	cd ./bin && zip get-investment-summary-by-symbol.zip bootstrap

	env GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -ldflags="-s -w" -o bin/bootstrap apps/investments/recalculate-avg-price/main.go
	cd ./bin && zip recalculate-avg-price.zip bootstrap

clean:
#	rm -rf ./bin ./vendor go.sum
	go clean
	rm -rf ./bin

deploy-prod: clean build
	npx sls deploy --stage prod --verbose

deploy-dev: clean build
	npx sls deploy --stage dev --verbose

remove-dev: clean build
	npx sls remove --stage dev --verbose

dev: clean build 
	npx serverless offline --useDocker --host 0.0.0.0

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

deploy-py-apps-dev:
	npx sls deploy --stage dev --config serverless.python.runtime.yml --verbose

deploy-py-apps-prod:
	npx sls deploy --stage prod --config serverless.python.runtime.yml --verbose
		