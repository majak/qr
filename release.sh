#!/bin/bash

mkdir -p public

GOOS=js GOARCH=wasm go build -o qart/main.wasm ./qart
cp qart/index.html public/
cp qart/wasm_exec.js public/
cp qart/main.wasm public/
cp qart/style.css public/
cp qart/android-chrome-192x192.png public/
cp qart/android-chrome-512x512.png public/
cp qart/favicon-16x16.png public/
cp qart/favicon-32x32.png public/
cp qart/favicon.ico public/

npx gh-pages -d public/