export DSN=postgres://pr_reviewer_user:pr_reviewer_password@localhost:5432/pr_reviewer_service?sslmode=disable

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING=${DSN}
export GOOSE_MIGRATION_DIR=./migrations

.PHONY: run
run:
	go run cmd/main.go

.PHONY: swag
swag:
	 swag init -g internal/infrastructure/http/controller/team.go

.PHONY: test-integration
test-integration:
	GIN_MODE=release go test -count=1 -v -run ^TestSuiteFunc$  ./tests/

.PHONY: goose-up
goose-up:
	goose up

.PHONY: compose-up
compose-up:
	docker-compose -p test up -d postgres

.PHONY: compose-rm
compose-rm:
	docker-compose -p test rm -fvs

.PHONY: up-and-test
up-and-test:
	make compose-up
	./scripts/wait_pg.sh
	make goose-up
	make test-integration
	make compose-rm
