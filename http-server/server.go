package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type countHandler struct {
	mu sync.Mutex
	n  int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "count is %d\n", h.n)
}

func main() {
	h1 := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #1!\n")
	}
	http.HandleFunc("/", h1)

	h2 := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc #2!\n")
	}
	http.HandleFunc("/endpoint", h2)

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}
	http.HandleFunc("/hello", helloHandler)

	http.Handle("/count", new(countHandler))

	// HTTP Server start
	log.Fatal(http.ListenAndServe(":8080", nil))
}
