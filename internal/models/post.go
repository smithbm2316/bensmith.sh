package models

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"bensmith.sh/internal/views"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func GeneratePosts() {
	// setup markdown parser
	var md = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
	)

	filepath.WalkDir(Site.ViewsDir, func(path string, d fs.DirEntry, err error) error {
		// handle errors with reading the directory
		if err != nil {
			return err
		}
		// skip processing directories and files that aren't markdown
		if d.IsDir() || path[len(path)-3:] != ".md" {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		var basename = fileInfo.Name()
		var filename = basename[:len(basename)-3]
		var slug = filepath.Join("words", filename)

		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Create the output directory.
		dir := filepath.Join(Site.BuildDir, slug)
		if err := os.MkdirAll(dir, 0755); err != nil && err != os.ErrExist {
			log.Fatalf("failed to create dir %q: %v", dir, err)
		}

		// Create the output file.
		name := filepath.Join(dir, "index.html")
		f, err := os.Create(name)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}

		// Convert the markdown to HTML, and pass it to the template.
		var buf bytes.Buffer
		if err := md.Convert([]byte(fileBytes), &buf); err != nil {
			log.Fatalf("failed to convert markdown to HTML: %v", err)
		}

		// Create an unsafe component containing raw HTML.
		content := Unsafe(buf.String())

		// Use templ to render the template containing the raw HTML.
		err = views.ContentPage(slug, content).Render(context.Background(), f)
		if err != nil {
			log.Fatalf("failed to write output file: %v", err)
		}

		return nil
	})
}
