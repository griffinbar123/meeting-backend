package main

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitializeRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/ws", EstablishWS)
	router.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(":8001", router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	InitializeRouter()
}
