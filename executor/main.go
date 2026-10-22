package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	boxes, err := strconv.Atoi(env("ISOLATE_BOXES", "8"))
	if err != nil || boxes < 1 {
		log.Fatal("ISOLATE_BOXES must be a positive integer")
	}
	runner := NewIsolate(env("ISOLATE_BIN", "isolate"), env("ISOLATE_CGROUPS", "") == "1", boxes)

	srv := &http.Server{
		Addr:              env("ADDR", ":8001"),
		Handler:           newServer(runner),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
	}
	log.Printf("executor listening on %s (%d boxes)", srv.Addr, boxes)
	log.Fatal(srv.ListenAndServe())
}
