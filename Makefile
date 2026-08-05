MODULES := modules/contracts modules/platform modules/chains/arc modules/mcp sdk/go
BIN_DIR := $(CURDIR)/bin

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	@echo "==> building modules/contracts"; (cd modules/contracts && go build ./...) || exit 1
	@echo "==> building modules/chains/arc"; (cd modules/chains/arc && go build ./...) || exit 1
	@echo "==> building sdk/go"; (cd sdk/go && go build ./...) || exit 1
	@echo "==> building modules/platform"; (cd modules/platform && go build -o $(BIN_DIR)/ ./...) || exit 1
	@echo "==> building modules/mcp"; (cd modules/mcp && go build -o $(BIN_DIR)/ ./...) || exit 1

.PHONY: test
test:
	@for m in $(MODULES); do \
		echo "==> testing $$m"; \
		(cd $$m && go test ./...) || exit 1; \
	done

.PHONY: vet
vet:
	@for m in $(MODULES); do \
		echo "==> vetting $$m"; \
		(cd $$m && go vet ./...) || exit 1; \
	done

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: tidy
tidy:
	@for m in $(MODULES); do \
		(cd $$m && go mod tidy); \
	done
	go work sync

.PHONY: run-gateway
run-gateway:
	cd modules/platform && go run ./cmd/gateway

.PHONY: run-mcp
run-mcp:
	cd modules/mcp && go run ./cmd/mcp-server
