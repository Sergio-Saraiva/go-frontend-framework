//go:build js && wasm

package timer

import (
	"fmt"
	"go-frontend-framework/component"
	"go-frontend-framework/router"
	"go-frontend-framework/signal"
	"time"
)

type Component struct {
	component.BaseComponent
	Time *signal.Signal[string]

	// Internal state to control the timer
	stopChan chan struct{}
}

func New() *Component {
	return &Component{
		Time: signal.New("Loading..."),
	}
}

// LIFECYCLE: Called when component appears
func (c *Component) OnInit() {
	fmt.Println("🕒 Timer Component Initialized")
	c.stopChan = make(chan struct{})

	// Start a background ticker
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				c.Time.Set(t.Format("15:04:05"))
				fmt.Println("Tick...")
			case <-c.stopChan:
				fmt.Println("🛑 Timer Stopped")
				return // Exit the goroutine
			}
		}
	}()
}

// LIFECYCLE: Called when component is removed
func (c *Component) OnDestroy() {
	fmt.Println("🗑️ Timer Component Destroyed")
	// Signal the goroutine to stop
	close(c.stopChan)
}

func (c *Component) GoHome() {
	router.Navigate("/")
}
