package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Do some startup work.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("startup code")
	defer log.Println("cleanup: closing db")
	defer log.Println("cleanup: flushing queues")
	defer log.Println("cleanup: stopping background jobs")

	// Define routes.
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	// Make the HTTP server.
	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	// Listen for requests in another goroutine.
	log.Println("server listening on", srv.Addr)
	go srv.ListenAndServe()

	// Create a channel where we will be notified of signals. Make sure it is
	// buffered or the signal package might drop the send.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Block until a signal is received.
	sig := <-signals

	log.Printf("%v signal received, shutting down", sig)
	srv.Shutdown(context.Background())
}

func handler(w http.ResponseWriter, r *http.Request) {
	id := time.Now().Nanosecond()
	log.Printf("request %d starting", id)
	time.Sleep(3 * time.Second)
	log.Printf("request %d done", id)
}
