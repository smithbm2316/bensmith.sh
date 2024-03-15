package views

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"bensmith.sh/models"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/xml"
)

// A data type that has all the data necessary to
// generate an Atom, RSS, and JSON feed
type Feed struct {
	Title       string
	Subtitle    string
	Language    string
	AbsoluteUrl string
	Name        string
	Email       string
	Posts       []*models.Post
}

// Create a new `Feed` struct which will have data necessary to
// generate an Atom, RSS, and JSON feed
func NewFeed(posts []*models.Post) Feed {
	return Feed{
		Title:       "Ben Smith’s Blog",
		Subtitle:    "Ben’s writings and thoughts about tech",
		Language:    "en",
		AbsoluteUrl: "https://bensmith.sh",
		Name:        "Ben Smith",
		Email:       "bsmithdev@mailbox.org",
		Posts:       posts,
	}
}

// Needed in Atom feed template to use the latest post's date as the content
// for the `<updated>` tag
func (feed Feed) GetNewestPostDate() string {
	return feed.Posts[0].FormatRFC3339()
}

// Generate and write a new Feed to our build directory
func (feed Feed) GenerateFeed(feedName string, dir string) {
	// load the text template
	templateName := fmt.Sprintf("%s.tmpl", feedName)
	tmpl := template.Must(
		template.ParseFiles(
			filepath.Join(Dirs.Views, "feeds", feedName),
		),
	)

	// create a new buffer and write the template to it
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, feed)
	if err != nil {
		log.Fatalf("failed to execute `%s` template: %v", templateName, err)
	}

	// create the file to write to
	outputDir := filepath.Join(Dirs.Build, dir)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}
	outputPath := filepath.Join(outputDir, feedName)
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed to create output file for `%s`: %v", outputPath, err)
	}
	defer file.Close()

	// minify the XML or JSON and write the buffer to the file
	var mimetype string
	m := minify.New()
	if feedName == "rss.json" {
		mimetype = "application/json"
		m.AddFunc(mimetype, json.Minify)
	} else {
		mimetype = "text/xml"
		m.AddFunc(mimetype, xml.Minify)
	}
	mw := m.Writer(mimetype, file)
	if mw.Write(buf.Bytes()); err != nil {
		log.Fatalf("Couldn't minify the `%s` feed, %v", feedName, err)
	}
	if err := mw.Close(); err != nil {
		log.Fatalf("Error executing the feed minfier's `io.Close` method, %v", err)
	}

	log.Printf("Created feed from `%s` at %s\n", feedName, outputPath)
}
