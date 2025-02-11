include .env

lint:
	golangci-lint run ./...  --config .golangci.pipeline.yaml

up:
	@docker-compose --file .docker/docker-compose.dev.yml up --build -d --remove-orphans

down:
	@docker-compose --file .docker/docker-compose.dev.yml down