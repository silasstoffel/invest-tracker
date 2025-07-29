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
