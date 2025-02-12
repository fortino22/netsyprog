package handlers

import (
	"encoding/json"
	"net/http"
)

type Student struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Subject string `json:"subject"`
}

var students = []Student{
	{ID: "1", Name: "John Doe", Age: 20, Subject: "Mathematics"},
	{ID: "2", Name: "Jane Smith", Age: 22, Subject: "Physics"},
}

func Route1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(students)
	}
}
