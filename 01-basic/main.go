package main

import (
	"log"
	"net/http"
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

	// Start listening for requests.
	log.Println("server listening on localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}

func handler(w http.ResponseWriter, r *http.Request) {
	id := time.Now().Nanosecond()
	log.Printf("request %d starting", id)
	time.Sleep(3 * time.Second)
	log.Printf("request %d done", id)
}
