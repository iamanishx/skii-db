package server

import (
	"encoding/json"
	"io"
	"net/http"
)

type RequestPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ResponsePayload struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func handlerSet(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var rp RequestPayload

		err = json.Unmarshal(body, &rp)
		if err != nil {
			http.Error(w, "Error parsing request payload", http.StatusBadRequest)
			return
		}

		err = e.Set(rp.Key, rp.Value)
		if err != nil {
			responseJSON(w, ResponsePayload{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusInternalServerError)
			return
		}

		responseJSON(w, ResponsePayload{
			Status:  "success",
			Message: "Key set successfully",
		}, http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}


func responseJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	d, err :=json.Marshal(data)
	w.Header().Set("content-type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(d)
}
