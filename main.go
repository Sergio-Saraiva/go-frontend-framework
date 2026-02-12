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
	router.Register("/", func() router.Component { return counter.New() })
	router.Register("/about", func() router.Component { return about.New() })
	root := router.Outlet()
	js.Global().Get("document").Get("body").Call("appendChild", root)
	router.Start()

	<-c
}
