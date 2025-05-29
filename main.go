package main

import (
	"fmt"
	"skii-db/engine"
)

func main() {
	e, err := engine.NewEngine()
	if err != nil {
		fmt.Println("Error initializing engine:", err)
		return
	}
	e.Set("foo", "bar")
	val, err := e.Get("key")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Retrieved value:", val)
}
