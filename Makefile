.PHONY: tidy run server docs docs-install fmt vet test clean build

GO          ?= go

tidy:
	$(GO) mod tidy

build:
	$(GO) build -o routeapi ./cmd/server

run server:
	$(GO) run ./cmd/server

# Install the swag CLI matching the version pinned in go.mod / tools.go.
docs-install:
	$(GO) install github.com/swaggo/swag/cmd/swag

# Regenerate docs/docs.go, docs/swagger.json, docs/swagger.yaml from the
# annotations on cmd/server/main.go and the handlers in internal/httpserver.
docs:
	@if [ -z "$(SWAG)" ] && ! command -v swag >/dev/null 2>&1; then \
	  echo "swag not found in PATH. Run: make docs-install"; \
	  exit 1; \
	fi
	swag init \
	  --generalInfo cmd/server/main.go \
	  --output docs \
	  --parseDependency \
	  --parseInternal

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

clean:
	rm -rf bin
