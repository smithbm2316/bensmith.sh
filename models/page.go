package models

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// The data processed from a markdown file in our `Dirs.Content` directory.
// Used to generate a specific `Page` on our site instead of having to author
// the page fully in HTML
type Page struct {
	Title string
	Slug  string
	// The generated HTML from a markdown source file
	Content string
}

// Instantiates a new `Page` object that will represent a specific route on our site
func NewPage(md goldmark.Markdown, metadataContext parser.Context, path string) (*Page, error) {
	rawMarkdownBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// create the slug from our markdown file's `path`
	slugWithFileExtension, ok := strings.CutPrefix(path, Dirs.Content)
	if !ok {
		log.Fatalf("failed to build a slug from %s", path)
	} else if slugWithFileExtension[:1] != "/" {
		slugWithFileExtension = "/" + slugWithFileExtension
	}
	slug, ok := strings.CutSuffix(slugWithFileExtension, ".md")
	if !ok {
		log.Fatalf("failed to build a slug from %s", path)
	}
	slug = slug + "/"

	// Convert the markdown to HTML, and pass it to the template.
	var markdownBuffer bytes.Buffer
	err = md.Convert(rawMarkdownBytes, &markdownBuffer, parser.WithContext(metadataContext))
	if err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	metadata := meta.Get(metadataContext)

	// Parse the title from the YAML frontmatter
	parseErrorMsg := "Couldn't parse `%s` from " + path
	title, ok := metadata["title"].(string)
	if !ok {
		log.Fatalf(parseErrorMsg, "title")
	}

	return &Page{
		Title:   title,
		Slug:    slug,
		Content: markdownBuffer.String(),
	}, err
}

// Convert the Page struct to a string representation
func (p Page) String() string {
	return fmt.Sprintf("%#v", p)
}
