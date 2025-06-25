package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    targetURL := os.Getenv("V2RAY_SERVER_URL") // Change to full URL
    if targetURL == "" {
        log.Fatal("V2RAY_SERVER_URL environment variable is required")
    }

    target, err := url.Parse(targetURL)
    if err != nil {
        log.Fatal("Invalid target URL:", err)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Create new request to target
        targetReq := r.Clone(r.Context())
        targetReq.URL.Scheme = target.Scheme
        targetReq.URL.Host = target.Host
        targetReq.RequestURI = ""

        client := &http.Client{}
        resp, err := client.Do(targetReq)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadGateway)
            return
        }
        defer resp.Body.Close()

        // Copy headers
        for key, values := range resp.Header {
            for _, value := range values {
                w.Header().Add(key, value)
            }
        }
        w.WriteHeader(resp.StatusCode)
        io.Copy(w, resp.Body)
    })

    log.Printf("HTTP proxy server listening on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}