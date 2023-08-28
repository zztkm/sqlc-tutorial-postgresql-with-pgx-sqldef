BIN := app
GOBIN ?= $(shell go env GOPATH)/bin

.PHONY: up
up:
	docker compose up -d --build --force-recreate

.PHONY: stop
stop:
	docker compose stop

.PHONY: down
down:
	docker compose down

.PHONY: logs-api
restart-api:
	docker compose restart api

.PHONY: logs-api
logs-api:
	docker compose logs -f -t api

.PHONY: ps
ps:
	docker compose ps


.PHONY: all
all: build

.PHONY: tag
tag:
	git tag "v${VERSION}"

.PHONY: build
build:
	go build -trimpath -o $(BIN) cmd/server/main.go

.PHONY: test
test: build
	go test -v ./...

.PHONY: lint
lint: $(GOBIN)/staticcheck
	staticcheck ./...

$(GOBIN)/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: gen
gen: $(GOBIN)/sqlc
	sqlc generate

$(GOBIN)/sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: migration
migration-dry-run: $(GOBIN)/psqldef
	psqldef -U test -W test -p 15432 testdb --dry-run < schema.sql

.PHONY: migration
migration: $(GOBIN)/psqldef
	psqldef -U test -W test -p 15432 testdb < schema.sql

$(GOBIN)/psqldef:
	go install github.com/k0kubun/sqldef/cmd/psqldef@latest

