package handlers

import (
 "fmt"
 "io"
 "net/http"
 "os"
 "path/filepath"
 "time"

 "github.com/Sevacoming/sprint6/internal/service"
)

func Index(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodGet {
  http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
  return
 }
 http.ServeFile(w, r, "index.html")
}

func Upload(w http.ResponseWriter, r *http.Request) {
 if r.Method != http.MethodPost {
  http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
  return
 }

 if err := r.ParseMultipartForm(10 << 20); err != nil {
  http.Error(w, "parse form error", http.StatusInternalServerError)
  return
 }

 f, fh, err := r.FormFile("file")
 if err != nil {
  http.Error(w, "get file error", http.StatusInternalServerError)
  return
 }
 defer f.Close()

 data, err := io.ReadAll(f)
 if err != nil {
  http.Error(w, "read file error", http.StatusInternalServerError)
  return
 }

 converted, err := service.DetectAndConvert(string(data))
 if err != nil {
  http.Error(w, "convert error: "+err.Error(), http.StatusInternalServerError)
  return
 }

 ext := filepath.Ext(fh.Filename)
 name := time.Now().UTC().Format("20060102T150405") + ext
 if err := os.WriteFile(name, []byte(converted), 0644); err != nil {
  http.Error(w, "save result error", http.StatusInternalServerError)
  return
 }

 w.Header().Set("Content-Type", "text/plain; charset=utf-8")
 fmt.Fprint(w, converted)
}
