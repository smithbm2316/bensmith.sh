package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"bensmith.sh/internal/models"
	"bensmith.sh/internal/views"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type Directories struct {
	Build   string
	Content string
	Posts   string
	Views   string
}

var Dirs = Directories{
	Build:   "src",
	Content: "internal/content",
	Posts:   "internal/content/words",
	Views:   "internal/views",
}

func main() {
	// True if we are in development mode
	// Whether or not the app is in development mode
	var DevMode bool
	flag.BoolVar(&DevMode, "dev", false, "True if we are in development mode")
	flag.Parse()
	// inject DevMode into "views" package so that we can include dev-mode only scripts and checks
	views.DevMode = DevMode

	// create output directory and generate posts
	if err := os.MkdirAll(Dirs.Build, os.ModePerm); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// setup markdown parser and metadata context
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, meta.Meta),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
	)
	metadataContext := parser.NewContext()

	// generate posts
	posts := GeneratePosts(md, metadataContext)

	// Create an index page
	indexBuildPath := filepath.Join(Dirs.Build, "index.html")
	indexFile, err := os.Create(indexBuildPath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	err = views.IndexPage().Render(context.Background(), indexFile)
	if err != nil {
		log.Fatalf("failed to write index page: %v", err)
	}
	fmt.Printf("Created %s at %s\n", "/", indexBuildPath)

	// And a blog index page
	blogIndexBuildPath := filepath.Join(Dirs.Build, "/words", "index.html")
	blogIndexFile, err := os.Create(blogIndexBuildPath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	// Write it out.
	err = views.BlogPage(posts).Render(context.Background(), blogIndexFile)
	if err != nil {
		log.Fatalf("failed to write blog index page: %v", err)
	}
	fmt.Printf("Created %s at %s\n", "/words", blogIndexBuildPath)

	fmt.Printf("Generated static files to %s\n", Dirs.Build)
}

func GeneratePosts(md goldmark.Markdown, metadataContext parser.Context) []*models.Post {
	var posts = make([]*models.Post, 0)

	filepath.WalkDir(Dirs.Posts, func(path string, d fs.DirEntry, err error) error {
		// handle errors with reading the directory
		if err != nil {
			return err
		}
		// skip processing directories and files that aren't markdown
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		post, err := models.NewPost(md, metadataContext, path)
		if err != nil {
			log.Fatalf("Failed to create new post from '%s'", path)
		}
		posts = append(posts, post)

		// Create the output directory.
		dir := filepath.Join(Dirs.Build, post.Slug)
		if err := os.MkdirAll(dir, 0755); err != nil && err != os.ErrExist {
			log.Fatalf("failed to create dir %q: %v", dir, err)
		}

		// Create the output file.
		outputPath := filepath.Join(dir, "index.html")
		f, err := os.Create(outputPath)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}

		// Use templ to render the template containing the raw HTML.
		err = views.PostRoute(post).Render(context.Background(), f)
		if err != nil {
			log.Fatalf("failed to write output file: %v", err)
		}

		fmt.Printf("Created %s at %s\n", post.Slug, outputPath)
		return nil
	})

	return posts
}
