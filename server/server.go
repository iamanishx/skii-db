package server

import (
	"fmt"
	"net/http"
	"skii-db/engine"
)

var e *engine.Engine


func Server() {
	var err error
	e, err = engine.NewEngine()
	if err != nil {
		panic("Error initializing engine: " + err.Error())
	}
	defer e.Close()
	e.Restore()

	go e.CompactFile()
	go e.DeleteFromFile()

	http.HandleFunc("/set", handlerSet)
	http.HandleFunc("/get", handlerGet)
	http.HandleFunc("/delete", handlerDelete)

	address := ":8080"
	fmt.Printf("Starting server on http://localhost%s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic("Error starting server: " + err.Error())
	}

}
