package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Sevacoming/sprint6/internal/handlers"
)

type Server struct {
	Logger *log.Logger
	HTTP   *http.Server
}

func New(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/upload", handlers.Upload)

	s := &http.Server{
		Addr:         ":8080", // строго 8080 для автотестов
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Server{Logger: logger, HTTP: s}
}

func (s *Server) Start() error {
	s.Logger.Println("server.go: server is listening on :8080")
	return s.HTTP.ListenAndServe()
}
