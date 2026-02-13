package component

import (
	"fmt"
	"go-frontend-framework/signal"
	"reflect"
)

type BaseComponent struct {
	target    interface{}
	Listeners map[string]func(any)
	Inputs    map[string]any
}

func (bc *BaseComponent) Init(target interface{}) {
	bc.target = target
}

func (bc *BaseComponent) SetInput(name string, value any) {
	if bc.target == nil {
		fmt.Println("Error: Component didn't call Base.Init(c) in New()")
		return
	}

	val := reflect.ValueOf(bc.target).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("input")
		if tag != name {
			continue
		}

		if field.CanInterface() {
			fieldVal := field.Interface()

			if setter, ok := fieldVal.(signal.AnySetter); ok {
				setter.SetAny(value)
				return
			}
		}

		if field.CanSet() {
			valToSet := reflect.ValueOf(value)

			if valToSet.Type().AssignableTo(field.Type()) {
				field.Set(valToSet)
				return
			}

			if valToSet.Type().ConvertibleTo(field.Type()) {
				field.Set(valToSet.Convert(field.Type()))
				return
			}

			fmt.Printf("Input Mismatch for '%s': Cannot assign %T to field %s (%s)\n",
				name, value, fieldType.Name, field.Type())
		}

		return
	}

	if bc.Inputs == nil {
		bc.Inputs = make(map[string]any)
	}

	bc.Inputs[name] = value
}

func (bc *BaseComponent) Emit(eventName string, payload any) {
	if bc.Listeners != nil {
		if callback, ok := bc.Listeners[eventName]; ok {
			callback(payload)
		}
	}
}

func (bc *BaseComponent) GetInput(name string) any {
	if bc.Inputs == nil {
		return nil
	}
	return bc.Inputs[name]
}

func (bc *BaseComponent) SetEventListener(event string, callback func(any)) {
	if bc.Listeners == nil {
		bc.Listeners = make(map[string]func(any))
	}
	bc.Listeners[event] = callback
}

func (bc *BaseComponent) OnInit()    {}
func (bc *BaseComponent) OnDestroy() {}
