include .env

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN=$(PG_MIGRATION_DSN)

lint:
	golangci-lint run ./...  --config .golangci.pipeline.yaml

up:
	@docker-compose --file .docker/docker-compose.dev.yml up --build -d --remove-orphans

down:
	@docker-compose --file .docker/docker-compose.dev.yml down


local-migration-status:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v