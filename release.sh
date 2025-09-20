#!/bin/bash

mkdir -p public

GOOS=js GOARCH=wasm go build -o qart/main.wasm ./qart
cp qart/index.html public/
cp qart/wasm_exec.js public/
cp qart/main.wasm public/

npx gh-pages -d public/