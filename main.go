//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/Sergio-Saraiva/go-frontend-framework/component"
	"github.com/Sergio-Saraiva/go-frontend-framework/di"
	"github.com/Sergio-Saraiva/go-frontend-framework/router"
	"github.com/Sergio-Saraiva/go-frontend-framework/src/app/counter"
	"github.com/Sergio-Saraiva/go-frontend-framework/src/app/timer"
	"github.com/Sergio-Saraiva/go-frontend-framework/src/services"
)

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Application Started")
	userService := services.NewUserService()
	di.Provide(userService)
	router.Register("/timer", func() component.Interface { return timer.New() })
	router.Register("/", func() component.Interface { return counter.New() })
	root := router.Outlet()
	js.Global().Get("document").Get("body").Call("appendChild", root)
	router.Start()

	<-c
}
