package httpserver

import (
	"log"
	"net/http"
)

func response(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(message)); err != nil {
		log.Printf("error writting HTTP response: %v", err)
	}
}
