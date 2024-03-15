package models

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/xml"
)

type SitemapRoute struct {
	Url          string
	LastModified time.Time
}

func (s SitemapRoute) FormatDate() string {
	return s.LastModified.Format(time.RFC3339)
}

type Sitemap struct {
	Routes []SitemapRoute
}

func NewSitemap() Sitemap {
	return Sitemap{
		Routes: make([]SitemapRoute, 0),
	}
}

// Generate and write a new Feed to our build directory
func (s Sitemap) Generate(slug string) {
	// load the sitemap.xml text template
	tmpl := template.Must(
		template.ParseFiles(
			filepath.Join(Dirs.Views, "sitemap.tmpl"),
		),
	)

	// create a new buffer and write the template to it
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, s.Routes)
	if err != nil {
		log.Fatalf("failed to execute sitemap.xml template: %v", err)
	}

	// create the file to write to
	file, err := os.Create(filepath.Join(Dirs.Build, slug))
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer file.Close()

	// minify the XML and write the buffer to the file
	m := minify.New()
	m.AddFunc("text/xml", xml.Minify)
	mw := m.Writer("text/xml", file)
	if mw.Write(buf.Bytes()); err != nil {
		log.Fatalf("Couldn't minify the XML sitemap, %v", err)
	}
	if err := mw.Close(); err != nil {
		log.Fatalf("Error executing the XML minfier's `io.Close` method, %v", err)
	}

	log.Printf("Created sitemap at %s\n", slug)
}
