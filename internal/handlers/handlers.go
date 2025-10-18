package handlers

import (
	"fmt"
	"net/http"

	"github.com/Sevacoming/sprint6/internal/service"
)

// GET / — отдать index.html
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

// POST /upload — принять файл под именем "file", автоопределить формат и вернуть результат
func UploadLegacy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// распарсить multipart (10 МБ буфер более чем)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file field 'file' required", http.StatusBadRequest)
		return
	}
	defer f.Close()

	converted, err := service.ConvertFromReader(f)
	if err != nil {
		http.Error(w, fmt.Sprintf("convert error: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(converted))
}
