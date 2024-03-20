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

// An entry representing one route in the sitemap
type SitemapEntry struct {
	Url          string
	LastModified time.Time
}

// Format the `SitemapEntry`'s `LastModified` date to ISO-8601/RFC3339.
// Used in the text/template that we generate our sitemap from
func (s SitemapEntry) FormatDate() string {
	return s.LastModified.Format(time.RFC3339)
}

// The data type used to generate the Sitemap for the site
type Sitemap struct {
	Entries []SitemapEntry
}

// Instantiates a new `Sitemap`
func NewSitemap() Sitemap {
	return Sitemap{
		Entries: make([]SitemapEntry, 0),
	}
}

// Generate and write a new Feed to our build directory
func (s Sitemap) Generate(slug string) {
	// load the sitemap.xml text template
	tmpl := template.Must(
		template.ParseFiles(
			filepath.Join(Dirs.TextTemplates, "sitemap.tmpl.xml"),
		),
	)

	// create a new buffer and write the template to it
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, s.Entries)
	if err != nil {
		log.Fatalf("failed to execute sitemap.xml template: %v", err)
	}

	// create the file to write to
	outputPath := filepath.Join(Dirs.Build, slug)
	file, err := os.Create(outputPath)
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

	log.Printf("Created %s\n", slug)
}
