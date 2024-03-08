package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"bensmith.sh/internal/models"
	"bensmith.sh/internal/views"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
)

type SiteSettings struct {
	BuildDir string
}

var Site = SiteSettings{
	BuildDir: "build",
}

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

	posts := models.GetPosts()

	// create output directory and generate posts
	if err := os.MkdirAll(Site.BuildDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}
	generatePosts(posts)

	// Create an index page.
	name := path.Join(Site.BuildDir, "index.html")
	f, err := os.Create(name)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	// Write it out.
	err = views.IndexPage(posts).Render(context.Background(), f)
	if err != nil {
		log.Fatalf("failed to write index page: %v", err)
	}

	// setup file server
	http.Handle("/", http.FileServer(http.Dir("./build")))

	log.Printf("Listening on :%d...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func generatePosts(posts []models.Post) {
	// Create a page for each post.
	for _, post := range posts {
		// Create the output directory.
		dir := path.Join(Site.BuildDir, post.Date.Format("2006/01/02"), post.Slug)
		if err := os.MkdirAll(dir, 0755); err != nil && err != os.ErrExist {
			log.Fatalf("failed to create dir %q: %v", dir, err)
		}

		// Create the output file.
		name := path.Join(dir, "index.html")
		f, err := os.Create(name)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}

		// Convert the markdown to HTML, and pass it to the template.
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(post.Content), &buf); err != nil {
			log.Fatalf("failed to convert markdown to HTML: %v", err)
		}

		// Create an unsafe component containing raw HTML.
		content := Unsafe(buf.String())

		// Use templ to render the template containing the raw HTML.
		err = views.ContentPage(post.Slug, content).Render(context.Background(), f)
		if err != nil {
			log.Fatalf("failed to write output file: %v", err)
		}
	}
}
