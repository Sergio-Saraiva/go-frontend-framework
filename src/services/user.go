package services

import "github.com/Sergio-Saraiva/go-frontend-framework/signal"

type UserService struct {
	Name       *signal.Signal[string]
	IsLoggedIn *signal.Signal[bool]
}

func NewUserService() *UserService {
	return &UserService{
		Name:       signal.New("Guest"),
		IsLoggedIn: signal.New(false),
	}
}

func (s *UserService) Login(name string) {
	s.Name.Set(name)
	s.IsLoggedIn.Set(true)
}

func (s *UserService) Logout() {
	s.Name.Set("Guest")
	s.IsLoggedIn.Set(false)
}
