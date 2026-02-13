//go:build js && wasm

package router

import "syscall/js"

type Initializer interface {
	Init(target any)
}

type Component interface {
	Render() js.Value
}

var routes = make(map[string]func() Component)
var outlet js.Value

func Register(path string, factory func() Component) {
	routes[path] = factory
}

func Outlet() js.Value {
	doc := js.Global().Get("document")
	outlet = doc.Call("createElement", "div")
	outlet.Set("id", "router-outlet")
	return outlet
}

func Navigate(path string) {
	js.Global().Get("history").Call("pushState", nil, "", path)
	renderPath(path)
}

func renderPath(path string) {
	if outlet.IsUndefined() || outlet.IsNull() {
		return
	}
	outlet.Set("innerHTML", "")
	factory, exists := routes[path]
	if !exists {
		return
	}

	comp := factory()

	if initer, ok := comp.(Initializer); ok {
		initer.Init(comp)
	}

	rootEl := comp.Render()
	outlet.Call("appendChild", rootEl)
}

func Start() {
	currentPath := js.Global().Get("location").Get("pathname").String()
	renderPath(currentPath)

	js.Global().Get("window").Call("addEventListener", "popstate", js.FuncOf(func(this js.Value, args []js.Value) any {
		newPath := js.Global().Get("location").Get("pathname").String()
		renderPath(newPath)
		return nil
	}))
}
