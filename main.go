//go:build js && wasm

package main

import (
	"fmt"
	"go-frontend-framework/src/app/counter"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("🚀 App Started")

	// 1. Create Component
	myCounter := counter.New()

	// 2. Render it (this method comes from counter_gen.go)
	rootElement := myCounter.Render()

	// 3. Mount to Body
	js.Global().Get("document").Get("body").Call("appendChild", rootElement)

	<-c
}
