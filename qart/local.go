// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// "go run qart/local.go" to get the live development environment on localhost:8080.

//go:build ignore

package main

import (
	"log"
	"net/http"
)

func main() {
	// This simple server serves all files from the "./qart" directory.
	// It is more robust for environments like WSL.
	fs := http.FileServer(http.Dir("./qart"))
	http.Handle("/", fs)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
