MODULES := modules/contracts modules/platform modules/chains/arc modules/mcp sdk/go

.PHONY: build
build:
	@for m in $(MODULES); do \
		echo "==> building $$m"; \
		(cd $$m && go build ./...) || exit 1; \
	done

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
