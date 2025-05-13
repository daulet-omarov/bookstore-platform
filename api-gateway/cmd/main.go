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
			proxyRequest("http://localhost:8000", w, r)

		case strings.HasPrefix(path, "/users"):
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
	// Create a new request with the same method, URL, and body
	req, err := http.NewRequest(r.Method, target+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	req.Header = r.Header

	// Send the request using http.DefaultClient
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers and status code
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Copy the response body
	io.Copy(w, resp.Body)
}
