.PHONY: test build clean fullclean localstack dev deploy

AWS_REGION ?= eu-west-3
AWS_PROFILE ?= serverless
LOCALSTACK_PORT ?= 4566
LOCAL_DYNAMODB_PORT ?= 8000

test:
	go test -coverprofile=coverage.out -cover -v ./...

build: test
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handletelegram HandleTelegramCommands/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getgospelandnotify GetGospelAndNotify/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/sendgospel SendGospel/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/ondemand OnDemand/main.go

clean:
	rm -rf ./bin ./vendor
	aws --endpoint-url=http://localhost:$(LOCALSTACK_PORT) sqs delete-queue --queue-url=http://localhost:$(LOCALSTACK_PORT)/000000000000/magnifibot 2>/dev/null || true
	aws --endpoint-url=http://localhost:$(LOCAL_DYNAMODB_PORT) dynamodb delete-table --table-name MagnifibotUser 2>/dev/null || true
	docker-compose down 2>/dev/null || true

fullclean: clean
	rm -rf ./docker
	docker-compose down -v
	colima stop

localstack:
	colima status 2>/dev/null || colima start
	mkdir -p ./docker/dynamodb
	docker-compose up -d && sleep 3
	aws --endpoint-url=http://localhost:$(LOCALSTACK_PORT) sqs create-queue --queue-name magnifibot 2>/dev/null || true
	aws --endpoint-url=http://localhost:$(LOCAL_DYNAMODB_PORT) dynamodb create-table \
		--table-name MagnifibotUser \
		--attribute-definitions AttributeName=ChatID,AttributeType=N \
		--key-schema AttributeName=ChatID,KeyType=HASH \
		--provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 2>/dev/null || true

dev: localstack
	go run main.go

deploy: clean build
	sls deploy -r $(AWS_REGION) --aws-profile $(AWS_PROFILE) --verbose
