//go:build js && wasm

package about

import "go-frontend-framework/router"

type Component struct{}

func New() *Component {
	return &Component{}
}

func (c *Component) GoHome() {
	router.Navigate("/")
}
