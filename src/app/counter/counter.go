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
	// 1. The List Signal
	Items *signal.Signal[[]string]
}

func New() *Component {
	return &Component{
		Count: signal.New(0),
		Show:  signal.New(true),
		// 2. Initialize with some data
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

// 3. Method to modify the list
func (c *Component) AddItem() {
	// Get current list
	currentList := c.Items.Get()

	// Create new item
	newItem := fmt.Sprintf("Fruit %d", len(currentList)+1)

	// Update signal (Triggers *gFor to re-render)
	// Note: We must create a NEW slice or append returns a new one
	c.Items.Set(append(currentList, newItem))
}
