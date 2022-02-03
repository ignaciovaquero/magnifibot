.PHONY: build clean fullclean localstack dev deploy

AWS_REGION ?= eu-west-3
AWS_PROFILE ?= serverless
LOCALSTACK_PORT ?= 4566

build:
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handletelegram HandleTelegramCommands/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getgospelandnotify GetGospelAndNotify/main.go

clean:
	rm -rf ./bin ./vendor
	aws --endpoint-url=http://localhost:$(LOCALSTACK_PORT) sqs delete-queue --queue-url=http://localhost:$(LOCALSTACK_PORT)/000000000000/magnifibot
	docker-compose down

fullclean: clean
	docker-compose down -v
	colima stop

localstack:
	colima status || colima start
	docker-compose up -d && sleep 3
	aws --endpoint-url=http://localhost:$(LOCALSTACK_PORT) sqs create-queue --queue-name magnifibot

dev: localstack
	go run main.go

deploy: clean build
	sls deploy -r $(AWS_REGION) --aws-profile $(AWS_PROFILE) --verbose
