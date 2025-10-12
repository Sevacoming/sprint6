package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Sevacoming/sprint6/internal/service"
)

// GET / — отдаем index.html
func IndexHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	}
}

// POST /upload — принимаем файл, конвертируем, сохраняем результат и возвращаем его
func UploadHandler(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(16 << 20); err != nil {
			logger.Printf("parse form error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			logger.Printf("form file error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			logger.Printf("read file error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		result, err := service.DetectAndConvert(string(data))
		if err != nil {
			logger.Printf("convert error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Имя результата: UTC timestamp + расширение исходного файла
		tstamp := time.Now().UTC().Format("20060102T150405Z0700")
		ext := filepath.Ext(header.Filename)
		if ext == "" {
			ext = ".txt"
		}
		outName := fmt.Sprintf("%s%s", tstamp, ext)

		if err := os.WriteFile(outName, []byte(result), 0o644); err != nil {
			logger.Printf("write result file error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(result))
	}
}
