package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case strings.HasPrefix(path, "/users") && strings.Contains(path, "/orders"):
			// Example: /users/123/orders/
			proxyRequest("http://localhost:8000", w, r)

		case strings.HasPrefix(path, "/users"):
			// Example: /users/, /users/123
			proxyRequest("http://localhost:4000", w, r)

		case strings.HasPrefix(path, "/books"):
			proxyRequest("http://localhost:4040", w, r)

		case strings.HasPrefix(path, "/orders"):
			proxyRequest("http://localhost:8000", w, r)

		default:
			http.NotFound(w, r)
		}
	})

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyRequest(target string, w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(target + r.URL.Path)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}
