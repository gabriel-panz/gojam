package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gabriel-panz/gojam/api"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

func init() {
	godotenv.Load()
}

func main() {
	c := cache.New(5*time.Minute, 10*time.Minute)
	server := api.NewServer(os.Getenv("ADDRESS"), os.Getenv("CLIENT_ID"), c)

	log.Printf("server now listening on %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
}
