.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -tags lambda.norpc -ldflags="-s -w" -o bin/bootstrap apps/investments/schedule/main.go
	cd ./bin && zip schedule-investment.zip bootstrap

clean:
#	rm -rf ./bin ./vendor go.sum
	go clean
	rm -rf ./bin

deploy: clean build
	npx sls deploy --stage dev --verbose

deploy-dev: clean build
	npx sls deploy --stage dev --verbose

remove-dev: clean build
	npx sls remove --stage dev --verbose	

dev: clean build 
	npx serverless offline --useDocker --host 0.0.0.0

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
