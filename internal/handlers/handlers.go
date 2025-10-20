package handlers

import (
	"io"
	"net/http"
	"os"
	"strings"

	service "github.com/Sevacoming/sprint6/internal/service"
)

// Index — отдаём корневой index.html из репозитория.
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Если файл есть — отдадим его, иначе вернём 500 (по замечанию ревьюера коды при ошибках)
	if _, err := os.Stat("index.html"); err == nil {
		http.ServeFile(w, r, "index.html")
		return
	}
	http.Error(w, "internal error", http.StatusInternalServerError)
}

// Upload читает multipart поле "myFile" (как в шаблоне), а также поддерживает:
// - fallback на поле "file" (для curl);
// - application/x-www-form-urlencoded: поле "text";
// - сырое text/plain тело.
// Коды: 200 — успех; 400 — пустой/непрочитанный ввод; 500 — ошибка конвертации.
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input string

	// 1) multipart: поле "myFile" (строго по шаблону)
	if f, _, err := r.FormFile("myFile"); err == nil {
		defer f.Close()
		if b, err := io.ReadAll(f); err == nil {
			input = string(b)
		}
	} else {
		// 1a) fallback: поле "file" (для совместимости с curl)
		if f2, _, err2 := r.FormFile("file"); err2 == nil {
			defer f2.Close()
			if b, err := io.ReadAll(f2); err == nil {
				input = string(b)
			}
		}
		// 2) x-www-form-urlencoded: поле "text"
		if input == "" {
			_ = r.ParseForm()
			if v := r.Form.Get("text"); v != "" {
				input = v
			}
		}
		// 3) сырое тело
		if input == "" {
			if b, err := io.ReadAll(r.Body); err == nil {
				_ = r.Body.Close()
				input = string(b)
			}
		}
	}

	input = strings.TrimSpace(input)
	if input == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	result, err := service.Convert(input)
	if err != nil {
		// по требованию ревьюера — 500 на ошибках конвертации
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, result)
}
