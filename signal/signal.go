package signal

var effectStack []func()
var currentEffect *EffectWrapper
var subIDCounter int

type EffectWrapper struct {
	run  func()
	deps []func()
}

type Signal[T comparable] struct {
	value       T
	subscribers map[int]func()
}

func New[T comparable](val T) *Signal[T] {
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
	if s.value == val {
		return
	}

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
