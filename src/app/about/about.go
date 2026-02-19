//go:build js && wasm

package about

import (
	"github.com/Sergio-Saraiva/go-frontend-framework/component"
	"github.com/Sergio-Saraiva/go-frontend-framework/router"
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
