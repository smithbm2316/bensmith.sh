package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// The port to run our server on
	var port int
	flag.IntVar(&port, "port", 2323, "The port to run our app's http server on")
	flag.Parse()

	// setup file server
	http.Handle("/", http.FileServer(http.Dir(".site")))

	log.Printf("Listening on :%d...", port)
	var err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
