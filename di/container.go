//go:build js && wasm

package di

import "reflect"

var registry = make(map[reflect.Type]any)

func Provide(service any) {
	t := reflect.TypeOf(service)
	registry[t] = service
}

func Get(t reflect.Type) any {
	return registry[t]
}

func Inject(target any) {
	val := reflect.ValueOf(target).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag.Get("inject")

		if tag == "true" && field.CanSet() {
			serviceType := field.Type()

			if service := Get(serviceType); service != nil {
				field.Set(reflect.ValueOf(service))
			}
		}
	}
}
