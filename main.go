//go:build js && wasm

package main

import (
	"fmt"
	"go-frontend-framework/component"
	"go-frontend-framework/di"
	"go-frontend-framework/router"
	"go-frontend-framework/src/app/counter"
	"go-frontend-framework/src/app/timer"
	"go-frontend-framework/src/services"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("🚀 Application Started")
	userService := services.NewUserService()
	di.Provide(userService)
	router.Register("/timer", func() component.Interface { return timer.New() })
	router.Register("/", func() component.Interface { return counter.New() })
	root := router.Outlet()
	js.Global().Get("document").Get("body").Call("appendChild", root)
	router.Start()

	<-c
}
