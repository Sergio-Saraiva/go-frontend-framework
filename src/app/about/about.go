//go:build js && wasm

package about

import (
	"go-frontend-framework/component"
	"go-frontend-framework/router"
)

type Component struct {
	component.BaseComponent
}

func New() *Component {
	return &Component{}
}

func (c *Component) GoHome() {
	router.Navigate("/")
}
