package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Route3Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		id := strings.TrimPrefix(r.URL.Path, "/route3/")

		var updatedStudent Student
		err := json.NewDecoder(r.Body).Decode(&updatedStudent)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		for i, student := range students {
			if student.ID == id {
				students[i] = updatedStudent
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(updatedStudent)
				return
			}
		}

		http.Error(w, "Student not found", http.StatusNotFound)
	}
}
