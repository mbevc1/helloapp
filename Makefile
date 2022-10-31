.PHONY: help
.DEFAULT: help
ifndef VERBOSE
.SILENT:
endif

NO_COLOR=\033[0m
GREEN=\033[32;01m
YELLOW=\033[33;01m
RED=\033[31;01m

VER?=dev
#GHASH:=$(shell git rev-parse --short HEAD)
#VERSION?=$(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo v0)
GO:=            go
GO_BUILD:=      go build -mod vendor -ldflags "-s -w -X main.GitCommit=${GHASH} -X main.Version=${VER}"
#VERSION="${VERSION}" goreleaser --snapshot --rm-dist
GO_VENDOR:=     go mod vendor
BIN:=           helloapp

help:: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-20s\033[0m %s\n", $$1, $$2}'

$(BIN): vendor ## Produce binary
	GO111MODULE=on $(GO_BUILD)

# We always want vendor to run **/*.go
.PHONY: vendor
vendor: *.go ## Build vendor deps
	GO111MODULE=on $(GO_VENDOR)

clean: clean-vendor ## Clean artefacts
	rm -rf $(BIN) $(BIN)_* $(BIN).exe go.sum
	go clean --modcache

clean-vendor: ## Clean vendor folder
	rm -rf vendor

build: clean $(BIN) ## Build binary
	upx $(BIN)

docker-build: ## Docker image build
	docker build -t $(BIN) .
	$(MAKE) docker-clean

docker-run: ## Docker image run
	docker run -it --rm --network host --name $(BIN) $(BIN)

docker-clean: ## Docker clean
	docker system prune -f

run:
	go run .

tidy:
	go mod tidy

dev: clean ## Dev test target
	go build -ldflags "-s -w -X main.Version=${VER}" -o $(BIN)
	upx $(BIN)

test: vendor ## Run tests
	go test -v ./...

fmt: *.go ## Formt Golang code
	go fmt ./...

lint:
	golint ./...

vet:
	go vet -all ./...

$(BIN)_linux_amd64: vendor **/*.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $@ *.go
	upx $@

$(BIN)_linux_alpine: vendor **/*.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $@ *.go
	upx $@

cache-start: ## start cache
	docker run -it --rm --network host --name cache redis

cache-stop: ## stop cache
	docker kill cache

db-start: ## start DB
	docker run -it --rm --network host --name db -e MYSQL_ROOT_PASSWORD=password1 mariadb

db-stop: ## stop DB
	docker kill db
