package signal

import (
	"fmt"
	"strconv"
)

var effectStack []func()
var currentEffect *EffectWrapper
var subIDCounter int

type EffectWrapper struct {
	run  func()
	deps []func()
}

type AnySetter interface {
	SetAny(val any)
}

type Signal[T any] struct {
	value       T
	subscribers map[int]func()
}

func New[T any](val T) *Signal[T] {
	return &Signal[T]{
		value:       val,
		subscribers: make(map[int]func()),
	}

}

func (s *Signal[T]) Subscribe(fn func()) func() {
	id := subIDCounter
	subIDCounter++

	s.subscribers[id] = fn

	return func() {
		delete(s.subscribers, id)
	}
}

func (s *Signal[T]) Get() T {
	if currentEffect != nil {
		cleanupFn := s.Subscribe(currentEffect.run)
		currentEffect.deps = append(currentEffect.deps, cleanupFn)
	}

	return s.value
}

func (s *Signal[T]) Set(val T) {
	// if s.value == val {
	// 	return
	// }

	s.value = val

	for _, fn := range s.subscribers {
		fn()
	}

}

func CreateEffect(fn func()) {
	e := &EffectWrapper{
		run: fn,
	}

	wrapper := func() {
		for _, unsubscribe := range e.deps {
			unsubscribe()
		}

		e.deps = nil
		prev := currentEffect
		currentEffect = e

		fn()

		currentEffect = prev
	}

	e.run = wrapper
	wrapper()
}

func (s *Signal[T]) SetAny(val any) {
	if casted, ok := val.(T); ok {
		s.Set(casted)
		return
	}

	var zero T
	switch any(zero).(type) {
	case int:
		if f, ok := val.(float64); ok {
			s.Set(any(int(f)).(T))
			return
		}

		if i, ok := val.(int); ok {
			s.Set(any(i).(T))
			return
		}

		if str, ok := val.(string); ok {
			if i, err := strconv.Atoi(str); err == nil {
				s.Set(any(i).(T))
			}
			return
		}

	case string:
		s.Set(any(fmt.Sprintf("%v", val)).(T))
		return

	case float64:
		if i, ok := val.(int); ok {
			s.Set(any(float64(i)).(T))
			return
		}

	case bool:
		if b, ok := val.(bool); ok {
			s.Set(any(b).(T))
			return
		}
	}

	fmt.Printf("Signal Type Mismatch: Signal[%T] cannot accept value of type %T (%v)\n", zero, val, val)

}
