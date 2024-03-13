package models

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type Post struct {
	Title        string
	Slug         string
	Tags         []string
	Content      string
	Published    time.Time
	LastModified time.Time
	Draft        bool
	Metadata     map[string]interface{}
}

func NewPost(md goldmark.Markdown, metadataContext parser.Context, path string) (*Post, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// create the slug from our markdown file's `path`
	slugWithFileExtension, ok := strings.CutPrefix(path, "internal/content")
	if !ok {
		log.Fatalf("failed to build a slug from %s", path)
	} else if slugWithFileExtension[:1] != "/" {
		slugWithFileExtension = "/" + slugWithFileExtension
	}
	slug, ok := strings.CutSuffix(slugWithFileExtension, ".md")
	if !ok {
		log.Fatalf("failed to build a slug from %s", path)
	}

	// Convert the markdown to HTML, and pass it to the template.
	var buf bytes.Buffer
	if err := md.Convert(fileBytes, &buf, parser.WithContext(metadataContext)); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	metadata := meta.Get(metadataContext)

	// Parse the metadata from the YAML frontmatter
	parseErrorMsg := "Couldn't parse `%s` from " + path
	title, ok := metadata["title"].(string)
	if !ok {
		log.Fatalf(parseErrorMsg, "title")
	}
	// these are optional fields
	tags, _ := metadata["tags"].([]string)
	draft, _ := metadata["draft"].(bool)

	// parse the `published` string date into a `time.Time` value
	published, ok := metadata["published"].(string)
	if !ok {
		log.Fatalf(parseErrorMsg, "published")
	}
	publishedDateParts := strings.Split(published, "-")
	year, err := strconv.Atoi(publishedDateParts[0])
	if err != nil {
		log.Fatalf("Couldn't parse the `published` date's `year` from '%s'", path)
	}
	month, err := strconv.Atoi(publishedDateParts[1])
	if err != nil {
		log.Fatalf("Couldn't parse the `published` date's `month` from '%s'", path)
	}
	day, err := strconv.Atoi(publishedDateParts[1])
	if err != nil {
		log.Fatalf("Couldn't parse the `published` date's `day` from '%s'", path)
	}

	return &Post{
		Title:        title,
		Slug:         slug,
		Tags:         tags,
		Content:      buf.String(),
		Published:    time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
		LastModified: time.Now().UTC(),
		Draft:        draft,
		Metadata:     metadata,
	}, err
}
