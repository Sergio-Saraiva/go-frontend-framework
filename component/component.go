//go:build js && wasm

package component

import "syscall/js"

type Interface interface {
	Render() js.Value
}

var registry = make(map[string]func() Interface)

func Register(tagName string, factory func() Interface) {
	registry[tagName] = factory
}

func Create(tagName string) Interface {
	factory, ok := registry[tagName]
	if !ok {
		return &ErrorComponent{Tag: tagName}
	}
	return factory()
}

type ErrorComponent struct {
	Tag string
}

func (e *ErrorComponent) Render() js.Value {
	doc := js.Global().Get("document")
	el := doc.Call("createElement", "div")
	el.Set("innerText", "⚠️ Component <"+e.Tag+"> not found")
	el.Get("style").Set("color", "red")
	return el
}
