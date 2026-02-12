package main

import (
	"fmt"
	"go-frontend-framework/cmd/compiler/signal"
)

func main() {
	count := signal.New(0)

	signal.CreateEffect(func() {
		fmt.Printf("The count is now %d\n", count.Get())
	})

	fmt.Println("Updating count...")
	count.Set(1)
	count.Set(5)
	count.Set(10)
}
