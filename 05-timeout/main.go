package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
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
	serverErrors := make(chan error, 1)
	go func() {
		log.Println("server listening on", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Create a channel where we will be notified of signals. Make sure it is
	// buffered or the signal package might drop the send.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Block until a signal is received or the server's listener fails.
	select {
	case err := <-serverErrors:
		if err != nil {
			return fmt.Errorf("server listening error: %v", err)
		}

	case sig := <-signals:
		log.Printf("%v signal received, shutting down", sig)

		// Define the maximum time you are willing to wait for HTTP requests to
		// finish. Keep in mind you have an overall deadline and there may be other
		// cleanup code waiting to finish.
		const gracePeriod = 2 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()

		// Start the shutdown process.
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("could not shut down in time: %v", err)
			return srv.Close()
		}
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	id := time.Now().Nanosecond()
	log.Printf("request %d starting", id)
	time.Sleep(3 * time.Second)
	log.Printf("request %d done", id)
}
