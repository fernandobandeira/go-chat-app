.PHONY: docker-compose-up build serve

docker-compose-up:
	docker-compose -f ./deploy/docker-compose.yaml up -d

build:
	go build -o go-chat-app *.go

serve: build
	chat serve
