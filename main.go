//go:build js && wasm

package main

import (
	"fmt"
	"go-frontend-framework/router"
	"go-frontend-framework/src/app/counter"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("🚀 Application Started")
	router.Register("/", func() router.Component { return counter.New() })
	root := router.Outlet()
	js.Global().Get("document").Get("body").Call("appendChild", root)
	router.Start()

	<-c
}
