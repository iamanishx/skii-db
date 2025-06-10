package server

import (
	"fmt"
	"net/http"
)

var e *Engine

func server() {
	var err error
	e, err = NewEngine()
	if err != nil {
		panic("Error initializing engine: " + err.Error())
	}
	defer e.Close()
	e.Restore()

	go e.CompactFile()
	go e.DeleteFile()

	http.HandleFunc("/set", handlerSet)
	http.HandleFunc("/get", handlerGet)
	http.HandleFunc("/delete", handlerDelete)

	address := ":8080"
	fmt.Println("Starting server on http://localhost%s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic("Error starting server: " + err.Error())
	}

}
