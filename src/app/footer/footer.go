//go:build js && wasm

package footer

import "fmt"

const Selector = "site-footer"

type Component struct{}

func New() *Component {
	return &Component{}
}

func (c *Component) Log() {
	fmt.Println("Footer button clicked!")
}
