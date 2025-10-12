package server

import (
 "log"
 "net/http"
 "os"
 "time"

 "github.com/Sevacoming/sprint6/internal/handlers"
)

type Server struct {
 Logger *log.Logger
 HTTP   *http.Server
}

func New(logger *log.Logger) *Server {
 mux := http.NewServeMux()
 mux.HandleFunc("/", handlers.IndexHandler(logger))
 mux.HandleFunc("/upload", handlers.UploadHandler(logger))

 port := os.Getenv("PORT")
 if port == "" {
  port = "8081"
 }

 httpSrv := &http.Server{
  Addr:         ":" + port,
  Handler:      mux,
  ErrorLog:     logger,
  ReadTimeout:  5 * time.Second,
  WriteTimeout: 10 * time.Second,
  IdleTimeout:  15 * time.Second,
 }

 return &Server{Logger: logger, HTTP: httpSrv}
}

func (s *Server) Start() error {
 s.Logger.Printf("server is listening on %s", s.HTTP.Addr)
 return s.HTTP.ListenAndServe()
}
