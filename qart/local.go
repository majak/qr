// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// "go run qart/local.go" to get the live development environment on localhost:8080.

//go:build ignore

package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// Proxy endpoint for the pixel art generator
	http.HandleFunc("/proxy-pixel-art", func(w http.ResponseWriter, r *http.Request) {
		// Forward the request to the Vercel server
		proxyReq, err := http.NewRequest(r.Method, "https://www.pixelart-pink.vercel.app/generate-pixel-art", r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		proxyReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy headers and status code from the Vercel server's response
		for name, values := range resp.Header {
			w.Header()[name] = values
		}
		w.WriteHeader(resp.StatusCode)

		// Copy the body from the Vercel server's response
		io.Copy(w, resp.Body)
	})

	// This simple server serves all files from the "./qart" directory.
	fs := http.FileServer(http.Dir("./qart"))
	http.Handle("/", fs)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
