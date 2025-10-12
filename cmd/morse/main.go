package main

import (
	"log"
	"os"

	"github.com/Sevacoming/sprint6/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "[morseapp] ", log.LstdFlags|log.Lshortfile)

	srv := server.New(logger)
	if err := srv.Start(); err != nil {
		logger.Fatal(err)
	}
}
