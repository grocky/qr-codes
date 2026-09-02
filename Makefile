VERSION := $(shell git describe --always --dirty 2>/dev/null || echo dev)

.DEFAULT_GOAL := help

.PHONY: help test e2e build build-wasm serve deploy tf-init tf-plan tf-apply tf-check tf-destroy clean

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*?##/ {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run Go unit tests
	go test ./...

build: ## Build all Go packages
	go build ./...

build-wasm: ## Build web/main.wasm and copy wasm_exec.js
	GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o web/main.wasm ./cmd/wasm
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/

serve: build-wasm ## Serve web/ locally on :8080
	@echo "serving web/ on http://localhost:8080"
	@cd web && python3 -m http.server 8080

e2e: build-wasm ## Run Playwright smoke tests
	cd web/e2e && npx playwright test

deploy: test build-wasm ## Deploy web/ to Cloudflare Pages
	npx wrangler pages deploy web --project-name $$(terraform -chdir=infrastructure output -raw pages_project)

tf-init: ## Initialize Terraform
	terraform -chdir=infrastructure init

tf-plan: ## Plan infrastructure changes
	terraform -chdir=infrastructure plan

tf-apply: ## Apply infrastructure changes
	terraform -chdir=infrastructure apply

tf-check: ## Validate and check formatting of Terraform files
	terraform -chdir=infrastructure validate
	terraform -chdir=infrastructure fmt -check -recursive

tf-destroy: ## Destroy infrastructure (careful!)
	terraform -chdir=infrastructure destroy

clean: ## Remove build artifacts
	rm -rf out web/main.wasm web/wasm_exec.js
