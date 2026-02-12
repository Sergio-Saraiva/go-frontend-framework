package counter

import "go-frontend-framework/signal"

type Component struct {
	// Capitalized so the generated file (same package) can see it
	Count *signal.Signal[int]
	Show  *signal.Signal[bool]
}

func New() *Component {
	return &Component{
		Count: signal.New(0),
		Show:  signal.New(true),
	}
}

func (c *Component) Increment() {
	c.Count.Set(c.Count.Get() + 1)
}

func (c *Component) Toggle() {
	c.Show.Set(!c.Show.Get())
}
