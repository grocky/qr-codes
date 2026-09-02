VERSION := $(shell git describe --always --dirty 2>/dev/null || echo dev)

.PHONY: test e2e build build-wasm serve deploy clean

test:
	go test ./...

build:
	go build ./...

build-wasm:
	GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o web/main.wasm ./cmd/wasm
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/

serve: build-wasm
	@echo "serving web/ on http://localhost:8080"
	@cd web && python3 -m http.server 8080

e2e: build-wasm
	cd web/e2e && npx playwright test

deploy: test build-wasm
	wrangler pages deploy web --project-name wifi-signs

clean:
	rm -rf out web/main.wasm web/wasm_exec.js
