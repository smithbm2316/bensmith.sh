package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"bensmith.sh/internal/models"
	"bensmith.sh/internal/views"
)

func main() {
	// True if we are in development mode
	// Whether or not the app is in development mode
	var DevMode bool
	flag.BoolVar(&DevMode, "dev", false, "True if we are in development mode")
	flag.Parse()
	// inject DevMode into "views" package so that we can include dev-mode only scripts and checks
	views.DevMode = DevMode

	// create output directory and generate posts
	if err := os.MkdirAll(models.Site.BuildDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}
	models.GeneratePosts()

	// Create an index page.
	name := path.Join(models.Site.BuildDir, "index.html")
	f, err := os.Create(name)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	// Write it out.
	err = views.IndexPage().Render(context.Background(), f)
	if err != nil {
		log.Fatalf("failed to write index page: %v", err)
	}

	fmt.Printf("Generated static files to %s\n", models.Site.BuildDir)
}
