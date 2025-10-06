package handlers

import "net/http"

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}
