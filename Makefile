.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/investments/create-investment apps/investments/create.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: deploy
	sls deploy --stage dev --verbose

deploy-dev: deploy
	sls deploy --stage dev --verbose

remove-dev: deploy
	sls remove --stage dev --verbose	

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
