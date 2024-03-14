package views

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"bensmith.sh/models"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/xml"
)

var Dirs models.Directories

type SitemapRoute struct {
	Url          string
	LastModified time.Time
}

func (s SitemapRoute) FormatDate() string {
	return s.LastModified.Format(time.RFC3339)
}

func GenerateSitemap(path string, sitemapRoutes []SitemapRoute) {
	// load the sitemap.xml text template
	tmpl := template.Must(
		template.ParseFiles(
			filepath.Join(Dirs.Views, "sitemap.tmpl"),
		),
	)

	// create a new buffer and write the template to it
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, sitemapRoutes)
	if err != nil {
		log.Fatalf("failed to execute sitemap.xml template: %v", err)
	}

	// create the file to write to
	file, err := os.Create(path)
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

	log.Printf("Created sitemap at %s\n", path)
}
