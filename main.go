package main

import (
    "fmt"
    "skii-db/engine"
)

func main() {
    e := engine.NewEngine()
    e.Set("key", "value")
    val, err := e.Get("key")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Retrieved value:", val)
}