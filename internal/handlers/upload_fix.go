package handlers

import (
	"io"
	"net/http"

	service "github.com/Sevacoming/sprint6/internal/service"
)

// Upload принимает multipart/form-data c полем "file".
// Возвращает:
//
//	200 text/plain — результат конвертации (без лишнего \n в конце);
//	400 — если вход некорректный или конвертация не удалась;
//	405 — для методов, отличных от POST.
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := service.ConvertFromReader(file)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, result) // важно: без завершающего \n
}
