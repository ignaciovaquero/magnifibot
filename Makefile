.PHONY: local-redis dev build clean

version = 0.1.0-SNAPSHOT

local-redis:
	COMPOSE_PROJECT_NAME=magnifibot docker-compose up -d

dev: redis
	go mod tidy
	go build -o bin/magnifibot
	./bin/magnifibot

build:
	docker build -t ivaquero/magnifibot:$(version)-amd64 --build-arg ARCH=amd64 . && \
	docker push ivaquero/magnifibot:$(version)-amd64 && \
	docker build -t ivaquero/magnifibot:$(version)-arm32v7 --build-arg ARCH=arm32v7 . && \
	docker push ivaquero/magnifibot:$(version)-arm32v7 && \
	docker build -t ivaquero/magnifibot:$(version)-arm64v8 --build-arg ARCH=arm64v8 . && \
	docker push	ivaquero/magnifibot:$(version)-arm64v8 && \
	docker manifest create \
		ivaquero/magnifibot:$(version) \
		ivaquero/magnifibot:$(version)-amd64 \
		ivaquero/magnifibot:$(version)-arm32v7 \
		ivaquero/magnifibot:$(version)-arm64v8 && \
	docker manifest push ivaquero/magnifibot:$(version)

clean:
	rm -rf ./bin
