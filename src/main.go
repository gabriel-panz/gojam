package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gabriel-panz/gojam/api"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	server := api.NewServer(os.Getenv("ADDRESS"), os.Getenv("CLIENT_ID"))

	log.Printf("server now listening on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
}
