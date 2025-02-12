package handlers

import (
	"encoding/json"
	"net/http"
)

func Route2Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newStudent Student
		err := json.NewDecoder(r.Body).Decode(&newStudent)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		students = append(students, newStudent)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newStudent)
	}
}
