package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Sevacoming/sprint6/internal/service"
)

// Index: отдаёт статический index.html
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

// Upload: принимает файл, определяет формат и конвертирует
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read file body: "+err.Error(), http.StatusBadRequest)
		return
	}

	out, err := service.DetectAndConvert(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// вернуть в ответ и записать в файл в корне проекта
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, out)

	_ = os.WriteFile("result.txt", []byte(out), 0644)
}
