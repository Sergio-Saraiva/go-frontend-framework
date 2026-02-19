//go:build js && wasm

package footer

import (
	"fmt"

	"github.com/Sergio-Saraiva/go-frontend-framework/component"
	"github.com/Sergio-Saraiva/go-frontend-framework/signal"
)

const Selector = "site-footer"

type Component struct {
	component.BaseComponent
	Status *signal.Signal[int] `input:"Status"`
}

func New() *Component {
	return &Component{
		Status: signal.New(0),
	}
}

func (c *Component) Log() {
	fmt.Println("Footer button clicked!")
}

func (c *Component) RequestReset() {
	c.Emit("reset", "Reset triggered from Child!")
}
