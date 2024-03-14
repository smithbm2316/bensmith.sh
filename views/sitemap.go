package views

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"time"
)

var sitemapTemplate = `<?xml version="1.0" encoding="utf-8"?><urlset xmlns="https://www.sitemaps.org/schemas/sitemap/0.9">{{ range . }}<url><loc>https://bensmith.sh{{ .Url }}</loc><lastmod>{{ .FormatDate }}</lastmod></url>{{ end }}</urlset>`

type SitemapRoute struct {
	Url          string
	LastModified time.Time
}

func (s SitemapRoute) FormatDate() string {
	return s.LastModified.Format(time.RFC3339)
}

func GenerateSitemap(path string, sitemapRoutes []SitemapRoute) {
	tmpl := template.Must(template.New("sitemap").Parse(sitemapTemplate))

	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, sitemapRoutes)
	if err != nil {
		log.Fatalf("failed to execute sitemap.xml template: %v", err)
	}

	fmt.Printf("Created sitemap at %s\n", path)
}
