include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=.$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=${DB_PORT} dbname=${DB_DB} user=${DB_USER} password=${DB_PASSWORD} sslmode=disable"

swagger:
	swag init -g cmd/main.go -o pkg/swagger --outputTypes json
	mv pkg/swagger/swagger.json pkg/swagger/api.swagger.json
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'


install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7

docker-build:
	docker compose up --build -d

docker-run:
	docker compose up -d

run:
	go run cmd/main.go

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v