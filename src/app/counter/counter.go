//go:build js && wasm

package counter

import (
	"fmt"
	"go-frontend-framework/router"
	"go-frontend-framework/signal"
)

type Component struct {
	Count *signal.Signal[int]
	Show  *signal.Signal[bool]
	Items *signal.Signal[[]string]
}

func New() *Component {
	return &Component{
		Count: signal.New(0),
		Show:  signal.New(true),
		Items: signal.New([]string{"Apple", "Banana", "Cherry"}),
	}
}

func (c *Component) GoToAbout() {
	router.Navigate("/about")
}

func (c *Component) Increment() {
	c.Count.Set(c.Count.Get() + 1)
}

func (c *Component) Toggle() {
	c.Show.Set(!c.Show.Get())
}

func (c *Component) AddItem() {
	currentList := c.Items.Get()
	newItem := fmt.Sprintf("Fruit %d", len(currentList)+1)
	c.Items.Set(append(currentList, newItem))
}
