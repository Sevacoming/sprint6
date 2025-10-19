package handlers

import (
	"io"
	"net/http"
	"strings"

	service "github.com/Sevacoming/sprint6/internal/service"
)

// Index — оставляем как простой health/page (если нужен по заданию)
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.WriteString(w, "OK")
}

// Upload принимает данные тремя способами (не ломая исходную форму):
// 1) multipart/form-data: поле "file"  (как в curl)
// 2) application/x-www-form-urlencoded: поле "text" (как на форме из задания)
// 3) другое: сырое тело запроса (text/plain; charset=utf-8)
// Успех: 200 text/plain; Ошибки конвертации: 500 (как просил ревьюер).
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input string

	// Попытка №1: multipart "file"
	if f, _, err := r.FormFile("file"); err == nil {
		defer f.Close()
		b, _ := io.ReadAll(f)
		input = string(b)
	} else {
		// Попытка №2: form field "text"
		_ = r.ParseForm()
		if v := r.Form.Get("text"); v != "" {
			input = v
		}
		// Попытка №3: сырое тело
		if input == "" {
			b, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()
			input = string(b)
		}
	}

	input = strings.TrimSpace(input)
	if input == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	result, err := service.Convert(input)
	if err != nil {
		// по замечанию ревьюера — 500 при ошибках
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, result) // без завершающего \n
}
