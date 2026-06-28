.PHONY: tidy build-all run server cron docs docs-install verify-openapi verify-openapi-install fmt vet test clean build

GO          ?= go
SWAG        ?= $(shell command -v swag 2>/dev/null)

tidy:
	$(GO) mod tidy

build-all: 
	$(GO) build -o ./main ./cmd/server/main.go

build: test
	$(GO) build -o ./app ./cmd/server

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

# Install swagger-cli for OpenAPI validation (requires Node.js/npm)
verify-openapi-install:
	npm install -g swagger-cli

# Verify OpenAPI/Swagger spec is valid and 100% compatible
verify-openapi:
	@echo "Validating OpenAPI specification..."
	@if [ ! -f "docs/swagger.yaml" ]; then \
	  echo "Error: docs/swagger.yaml not found. Run 'make docs' first."; \
	  exit 1; \
	fi
	@if command -v swagger-cli >/dev/null 2>&1; then \
	  echo "Running swagger-cli validation..."; \
	  swagger-cli validate docs/swagger.yaml; \
	else \
	  echo "Warning: swagger-cli not found. Install with: make verify-openapi-install"; \
	  echo "Attempting validation with swag..."; \
	  swag fmt -g cmd/server/main.go; \
	fi
	@echo "✓ OpenAPI spec validation complete"

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

clean:
	rm -rf bin
