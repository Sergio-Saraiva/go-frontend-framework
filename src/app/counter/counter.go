//go:build js && wasm

package counter

import (
	"fmt"

	"github.com/Sergio-Saraiva/go-frontend-framework/component"
	"github.com/Sergio-Saraiva/go-frontend-framework/router"
	"github.com/Sergio-Saraiva/go-frontend-framework/signal"
	"github.com/Sergio-Saraiva/go-frontend-framework/src/services"
)

// Default selector "counter" (derived from folder name)
type Component struct {
	component.BaseComponent

	Count *signal.Signal[int]

	User *services.UserService `inject:"true"`
}

func New() *Component {
	return &Component{
		Count: signal.New(0),
	}
}

func (c *Component) Increment() {
	c.Count.Set(c.Count.Get() + 1)
}

// OnReset is the Event Handler for the child's (reset) event.
func (c *Component) OnReset(payload any) {
	fmt.Println("Parent received reset event:", payload)

	// Reset the state to 0.
	// This will automatically propagate down to the child via [Status]!
	c.Count.Set(0)
}

func (c *Component) GoToTimer() {
	router.Navigate("/timer")
}
