package main

import (
	"fmt"
	"net/http"
	"time"
	"var2/handlers"
	"var2/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/route1", middleware.MethodValidation(handlers.Route1Handler))
	mux.HandleFunc("/route2", middleware.MethodValidation(handlers.Route2Handler))
	mux.HandleFunc("/route3/", middleware.MethodValidation(handlers.Route3Handler))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Println("Server is running on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}
