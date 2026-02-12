//go:build js && wasm

package main

import (
	"fmt"
	"go-frontend-framework/router"
	"go-frontend-framework/src/app/about"
	"go-frontend-framework/src/app/counter"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("🚀 App Started")

	// 1. Register Routes
	// We wrap the "New" function to match the func() Component signature
	router.Register("/", func() router.Component { return counter.New() })
	router.Register("/about", func() router.Component { return about.New() })

	// 2. Create the Outlet (The placeholder)
	// <div id="router-outlet"></div>
	root := router.Outlet()
	js.Global().Get("document").Get("body").Call("appendChild", root)

	// 3. Start the Router
	// This renders the component matching the current URL
	router.Start()

	<-c
}
