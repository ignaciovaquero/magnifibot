.PHONY: build clean dev deploy

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handletelegram HandleTelegramCommands/main.go

clean:
	rm -rf ./bin ./vendor go.sum

dev:
	go run main.go

deploy: clean build
	sls deploy --verbose
