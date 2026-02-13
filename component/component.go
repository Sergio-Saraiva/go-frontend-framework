//go:build js && wasm

package component

import "syscall/js"

type Interface interface {
	Render() js.Value
	SetInput(name string, value any)
	SetEventListener(eventName string, handler func(any))
	Init(target any)
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
	comp := factory()

	comp.Init(comp)

	return comp
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

func (e *ErrorComponent) Init(target any)                                   {}
func (e *ErrorComponent) SetInput(name string, value any)                   {}
func (e *ErrorComponent) SetEventListener(event string, callback func(any)) {}
