include .env

lint:
	golangci-lint run ./...  --config .golangci.pipeline.yaml
