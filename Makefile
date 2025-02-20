include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=.$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=${DB_PORT} dbname=${DB_DB} user=${DB_USER} password=${DB_PASSWORD} sslmode=disable"

swagger:
	swag init -g cmd/main.go -o pkg/swagger --outputTypes json
	mv pkg/swagger/swagger.json pkg/swagger/api.swagger.json
	$(LOCAL_BIN)/statik -src=pkg/swagger/ -include='*.css,*.html,*.js,*.json,*.png'

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5

lint:
	go mod tidy
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

docker-build:
	docker compose up -d --build

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

test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/merynayr/AvitoShop/internal/service/... -count 5

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/merynayr/AvitoShop/internal/service/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore


TEST_MIGRATION_DSN="host=localhost port=5433 dbname=avito_test user=test_user password=test_password sslmode=disable"

migration_for_test_up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${TEST_MIGRATION_DSN} up -v

migration_for_test_down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${TEST_MIGRATION_DSN} down -v

load-testing:
	k6 run internal/tests/load_test.js