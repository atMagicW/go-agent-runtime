APP_NAME=go-agent-runtime
CONFIG=configs/app.yaml

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...

.PHONY: docker-up
docker-up:
	docker compose -f deployments/compose/docker-compose.yaml up -d

.PHONY: docker-down
docker-down:
	docker compose -f deployments/compose/docker-compose.yaml down

.PHONY: migrate-print
migrate-print:
	@echo "Run SQL files in migrations/ manually or with your migration tool."

.PHONY: curl-health
curl-health:
	curl localhost:8080/v1/health

.PHONY: curl-capabilities
curl-capabilities:
	curl localhost:8080/v1/capabilities

.PHONY: test-short
test-short:
	go test ./internal/... ./api/...

.PHONY: check
check: fmt test