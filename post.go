package bs

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
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// Represents a heading in our `Post`
type Heading struct {
	Text string
	// represents the corresponding `<h2/3/4/5/6>` element
	Level int
	// the HTML id that we will link to in our table of contents on a `Post` page
	Id string
}

// The data processed from a markdown file in our `Dirs.Posts` directory.
// Used to generate a specific `Post` on our site instead of having to author
// the blog post fully in HTML
type Post struct {
	Title        string
	Slug         string
	Published    time.Time
	LastModified time.Time
	Headings     []Heading
	Tags         []string
	// a generic reference to all the YAML metadata parsed from the markdown file
	Metadata map[string]interface{}
	// The generated HTML from a markdown source file
	Content string
	// set to true to hide this post in production
	Draft bool
}

// Instantiates a new `Post` for the supplied `path`. Returns a pointer
// to that `Post` or an error
func NewPost(md goldmark.Markdown, metadataContext parser.Context, path string) (*Post, error) {
	// read the whole file into a byte slice in memory
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
	re := regexp.MustCompile(`(?m)^(?P<level>#{1,6})\s(?P<heading>.*)\s*?{(?P<id>#.*)}\s*?$`)
	matchNames := re.SubexpNames()
	headings := []Heading{}
	for _, match := range re.FindAllStringSubmatch(string(rawMarkdownBytes), -1) {
		nextHeading := Heading{}
		for i, matchedValue := range match {
			// if our capture group name is "level", we want to record the number of "#"s that we matched
			// if our capture group name is "heading", we want to save the raw text of the heading
			// if our capture group name is "id", we want to save the text of the heading's html ID
			switch matchNames[i] {
			case "level":
				nextHeading.Level = len(matchedValue)
			case "heading":
				nextHeading.Text = matchedValue
			case "id":
				nextHeading.Id = matchedValue
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

// Format the `Post`'s `Published` date to ISO-8601/RFC3339.
// Used in the text/template we generate a `Feed` from
func (p Post) FormatRFC3339() string {
	return p.Published.Format(time.RFC3339)
}

// Format the `Post`'s `Published` date to RFC822.
// Used in the text/template we generate a rss `Feed` from
func (p Post) FormatRFC822() string {
	return p.Published.Format(time.RFC822)
}

// Format the `Post`'s `Slug` to a URL that is prefixed
// with our domain name. Used in the text/template we generate
// a `Feed` from
func (p Post) FormatAbsoluteUrl() string {
	return fmt.Sprintf("https://bensmith.sh%s", p.Slug)
}
