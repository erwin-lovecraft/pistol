PROJECT_NAME := pistol
TEAM_NAME := viebiz
ENV := dev

# Run cmd
.PHONY: run
run:
	@go run ./cmd/serverd

.PHONY: sqlc
sqlc:
	@docker run --rm -v .:/app -w /app sqlc/sqlc generate
