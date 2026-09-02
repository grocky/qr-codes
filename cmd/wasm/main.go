//go:build js && wasm

// Command wasm exposes internal/sign to the browser as globalThis.wifiSign.
// It is glue only: all logic lives (and is tested) in internal/sign.
package main

import (
	"syscall/js"

	"qr-codes/internal/sign"
)

var version = "dev" // overridden at build time via -ldflags "-X main.version=..."

func main() {
	api := js.Global().Get("Object").New()
	api.Set("render", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 || args[0].Type() != js.TypeString {
			return `{"errors":[{"field":"","message":"render expects one JSON string argument"}]}`
		}
		return string(sign.Render([]byte(args[0].String())))
	}))
	api.Set("version", version)
	js.Global().Set("wifiSign", api)

	// Signal readiness to the page, then keep the Go runtime alive.
	if ready := js.Global().Get("wifiSignReady"); ready.Type() == js.TypeFunction {
		ready.Invoke()
	}
	select {}
}
