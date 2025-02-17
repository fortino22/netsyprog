package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Route4Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) 
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing multipart form: %v", err), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files received", http.StatusBadRequest)
		return
	}

	uploadDir := "../upload" 

	var uploadedFiles []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filePath := filepath.Join(uploadDir, fileHeader.Filename)
		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to write file: %v", err), http.StatusInternalServerError)
			return
		}

		uploadedFiles = append(uploadedFiles, fileHeader.Filename)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Files uploaded successfully: %v", uploadedFiles)
}
