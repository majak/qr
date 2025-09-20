// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// "go run qart/local.go" to get the live development environment on localhost:8080.

//go:build ignore

package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func main() {
	// Proxy endpoint for the pixel art generator
	http.HandleFunc("/proxy-pixel-art", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Proxy: Received request for /proxy-pixel-art")

		// Read the body from the incoming request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Proxy Error: reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		r.Body.Close()

		// Forward the request to the Vercel server
		proxyReq, err := http.NewRequest(r.Method, "https://pixelart-pink.vercel.app/generate-pixel-art", bytes.NewReader(body))
		if err != nil {
			log.Printf("Proxy Error: creating new request: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		proxyReq.Header.Set("Content-Type", "application/json")

		log.Println("Proxy: Forwarding request to Vercel...")
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			log.Printf("Proxy Error: contacting Vercel server: %v", err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		log.Printf("Proxy: Received response from Vercel. Status: %s", resp.Status)

		// Copy headers from the Vercel server's response
		for name, values := range resp.Header {
			w.Header()[name] = values
		}

		// Read the response body so we can log it and send it to the client
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Proxy Error: reading Vercel response body: %v", err)
			http.Error(w, "Error reading upstream response", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Proxy: Vercel returned error. Body: %s", string(respBody))
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	})

	// This simple server serves all files from the "./qart" directory.
	fs := http.FileServer(http.Dir("./qart"))
	http.Handle("/", fs)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
