# README

## About

This is the official Wails Svelte-TS template.

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

## Adding new suspect

1. place the suspect images into `./frontend/public/input` directory
2. run `go run dev/main.go import` - this will iterate over all images in `input` directory and parse them
3. take parsed images from `./frontend/public/input` and place them to `.frontend/public/suspects`
4. filename of the image is its sha256 sum

Sha256 sums are used to prevent duplicates.
Also UUID is taken from the sha256 filename (and not generated) to ensure Suspects have same UUIDs across all game instances.
This way statistics can be made from multiple computers.
