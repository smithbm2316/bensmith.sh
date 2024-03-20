package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"bensmith.sh/models"
)

var Dirs = models.Directories{
	Build: "www",
}

type fileSystemHandlingNotFound struct {
	root http.Dir
}

func (fs fileSystemHandlingNotFound) Open(name string) (http.File, error) {
	file, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		// File not found, serve the custom 404.html
		return fs.root.Open("/404.html")
	}
	return file, err
}

func main() {
	// The port to run our server on
	var port int
	flag.IntVar(&port, "port", 2323, "The port to run our app's http server on")
	flag.Parse()

	// Create a custom file server with our custom Open method and use it
	fileServer := fileSystemHandlingNotFound{root: http.Dir(Dirs.Build)}
	http.Handle("/", http.FileServer(fileServer))

	log.Printf("Listening on :%d...", port)
	var err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
