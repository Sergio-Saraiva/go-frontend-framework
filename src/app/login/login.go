//go:build js && wasm

package login

import (
	"github.com/Sergio-Saraiva/go-frontend-framework/component"
	"github.com/Sergio-Saraiva/go-frontend-framework/signal"
	"github.com/Sergio-Saraiva/go-frontend-framework/src/services"
)

const Selector = "app-login"

type Component struct {
	component.BaseComponent

	Name *signal.Signal[string]
	// Inject the SAME service instance
	User *services.UserService `inject:"true"`
}

func New() *Component {
	return &Component{
		Name: signal.New(""),
	}
}

func (c *Component) DoLogin() {
	c.User.Login(c.Name.Get()) // Updates global state
}
