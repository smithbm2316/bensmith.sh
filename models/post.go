package models

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type Heading struct {
	Level int
	Text  string
}

type Post struct {
	Title        string
	Slug         string
	Tags         []string
	Content      string
	Published    time.Time
	LastModified time.Time
	Draft        bool
	Metadata     map[string]interface{}
	Headings     []Heading
}

func NewPost(md goldmark.Markdown, metadataContext parser.Context, path string) (*Post, error) {
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

	// Parse the metadata from the YAML frontmatter
	parseErrorMsg := "Couldn't parse `%s` from " + path
	title, ok := metadata["title"].(string)
	if !ok {
		log.Fatalf(parseErrorMsg, "title")
	}

	// tags and draft are optional fields
	var tags []string
	tagsInterface := metadata["tags"]
	switch tagsStr := tagsInterface.(type) {
	case string:
		for _, tag := range strings.Split(tagsStr, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	case []string:
		for _, tag := range tagsStr {
			tags = append(tags, strings.TrimSpace(tag))
		}
	case []interface{}:
		for _, tag := range tagsStr {
			tags = append(tags, strings.TrimSpace(tag.(string)))
		}
	default:
		break
	}

	draft, _ := metadata["draft"].(bool)

	// parse the headings so we can build a TOC later
	re := regexp.MustCompile(`(?m)^(?P<level>#{1,6})\s(?P<heading>.*)$`)
	matchNames := re.SubexpNames()
	headings := []Heading{}
	for _, match := range re.FindAllStringSubmatch(string(rawMarkdownBytes), -1) {
		nextHeading := Heading{}
		for i, matchedValue := range match {
			// if our capture group name is "level", we want to record the number of "#"s that we matched
			// if our capture group name is "heading", we want to save the raw text of the heading
			switch matchNames[i] {
			case "level":
				nextHeading.Level = len(matchedValue)
			case "heading":
				nextHeading.Text = matchedValue
			}
		}
		headings = append(headings, nextHeading)
	}

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
		Content:      markdownBuffer.String(),
		Published:    time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
		LastModified: time.Now().UTC(),
		Draft:        draft,
		Metadata:     metadata,
		Headings:     headings,
	}, err
}

// Convert the Post struct to a string representation
func (p Post) String() string {
	return fmt.Sprintf("%#v", p)
}

func (p Post) FormatRFC3339() string {
	return p.Published.Format(time.RFC3339)
}

func (p Post) FormatRFC822() string {
	return p.Published.Format(time.RFC822)
}

func (p Post) FormatAbsoluteUrl() string {
	return fmt.Sprintf("https://bensmith.sh%s", p.Slug)
}
