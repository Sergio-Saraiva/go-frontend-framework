//go:build js && wasm

package router

import (
	"fmt"
	"syscall/js"

	"github.com/Sergio-Saraiva/go-frontend-framework/component"
)

var routes = make(map[string]func() component.Interface)
var outlet js.Value
var currentComp component.Interface

func Register(path string, factory func() component.Interface) {
	routes[path] = factory
}

func Outlet() js.Value {
	doc := js.Global().Get("document")
	outlet = doc.Call("createElement", "div")
	outlet.Set("id", "router-outlet")
	return outlet
}

func Navigate(path string) {
	win := js.Global().Get("window")
	win.Get("location").Set("hash", "#"+path)
}

func renderPath(path string) {
	fmt.Println("--- Router: Navigating to", path, "---")

	if currentComp != nil {
		currentComp.OnDestroy()
	}

	outlet.Set("innerHTML", "")
	factory, exists := routes[path]
	if !exists {
		outlet.Set("innerText", "404 - Not Found")
		return
	}

	comp := factory()

	currentComp = comp

	comp.Init(comp)
	comp.OnInit()

	rootEl := comp.Render()
	outlet.Call("appendChild", rootEl)
}

func Start() {
	win := js.Global().Get("window")

	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		hash := win.Get("location").Get("hash").String()
		path := "/"
		if len(hash) > 1 {
			path = hash[1:]
		}
		renderPath(path)
		return nil
	})
	win.Call("addEventListener", "hashchange", cb)

	hash := win.Get("location").Get("hash").String()
	path := "/"
	if len(hash) > 1 {
		path = hash[1:]
	}
	renderPath(path)
}
