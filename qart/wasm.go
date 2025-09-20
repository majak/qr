// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

// Qart is a WebAssembly program to help create artistic QR code images.
// The algorithms are described at https://research.swtch.com/qart,
// and this program is running at https://research.swtch.com/qr/draw/.
//
// To run the program locally, use “go run local.go”.
package main

import (
	"encoding/base64"
	"fmt"
	"syscall/js"
)

// generateQRCode is a stateless function that can be called from JavaScript.
// It takes image data and a URL and returns a base64-encoded PNG of the QR code.
func generateQRCode(this js.Value, args []js.Value) any {
	imageData := args[0]
	url := args[1].String()

	// Copy the image data from JavaScript to Go
	data := make([]byte, imageData.Get("length").Int())
	js.CopyBytesToGo(data, imageData)

	// Create a new Image object for this specific request
	pic := &Image{
		File:    data,
		URL:     url,
		Version: 6, // Default version
		Mask:    2, // Default mask
	}

	// Generate the QR code
	pngData, err := pic.Encode()
	if err != nil {
		// In case of an error, we can return a descriptive string.
		// A more robust solution would be to return a structured error object.
		return fmt.Sprintf("Error: %v", err)
	}

	// Return the base64-encoded PNG data
	return base64.StdEncoding.EncodeToString(pngData)
}

func main() {
	// Export the stateless function to be called from JavaScript
	js.Global().Set("generateQRCode", js.FuncOf(generateQRCode))

	// Keep the Go program running
	<-make(chan bool)
}
