package main

import (
	"log"
	"net/http"
	"time"
)

// Starter code for the graceful shutdown talk

func main() {

	// Do some startup work.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("startup code")

}

func handler(w http.ResponseWriter, r *http.Request) {
	id := time.Now().Nanosecond()
	log.Printf("request %d starting", id)
	time.Sleep(3 * time.Second)
	log.Printf("request %d done", id)
}
