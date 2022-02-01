.PHONY: build clean dev deploy

AWS_REGION ?= eu-west-3

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handletelegram HandleTelegramCommands/main.go

clean:
	rm -rf ./bin ./vendor

dev:
	go run main.go

deploy: clean build
	sls deploy -r $(AWS_REGION) --verbose

