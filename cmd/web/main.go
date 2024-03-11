package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"bensmith.sh/internal/views"
)

func main() {
	// True if we are in development mode
	// Whether or not the app is in development mode
	var DevMode bool
	flag.BoolVar(&DevMode, "dev", false, "True if we are in development mode")
	// The port to run our server on
	var port int
	flag.IntVar(&port, "port", 2323, "The port to run our app's http server on")
	flag.Parse()
	// inject DevMode into "views" package so that we can include dev-mode only scripts and checks
	views.DevMode = DevMode

	// setup file server
	http.Handle("/", http.FileServer(http.Dir(".site")))

	log.Printf("Listening on :%d...", port)
	var err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
