package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"bensmith.sh/models"
	"bensmith.sh/views"

	chroma "github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var Dirs = models.Directories{
	Build:   "src",
	Content: "content",
	Posts:   "content/words",
	Views:   "views",
}

func main() {
	// inject Dirs into "models" package
	models.Dirs = Dirs

	// create build directory for output
	if err := os.MkdirAll(Dirs.Build, os.ModePerm); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// setup file writer for chroma CSS generation, we want to write it to a static CSS file that targets the "chroma" class instead of using inline HTML styles
	chromaCSSPath := filepath.Join(Dirs.Build, "chroma.css")
	chromaCSSFile, err := os.Create(chromaCSSPath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer chromaCSSFile.Close()

	// setup markdown parser and metadata context
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle("catppuccin-mocha"),
				// consume the io.Write we just created
				highlighting.WithCSSWriter(chromaCSSFile),
				highlighting.WithFormatOptions(
					chromahtml.TabWidth(2),
					// make sure we use class generation over inline styles
					chromahtml.WithClasses(true),
					// make sure that <pre> tags overflow
					chromahtml.WithCustomCSS(map[chroma.TokenType]string{
						chroma.PreWrapper: "overflow-x: auto; padding: 1em; border-radius",
					}),
				),
			),
		),
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
	err = views.IndexRoute().Render(context.Background(), indexFile)
	if err != nil {
		log.Fatalf("failed to write index page: %v", err)
	}
	indexFile.Close()
	fmt.Printf("Created %s at %s\n", "/", indexBuildPath)

	// And a blog index page
	blogIndexBuildPath := filepath.Join(Dirs.Build, "/words", "index.html")
	blogIndexFile, err := os.Create(blogIndexBuildPath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	// Write it out.
	err = views.BlogRoute(posts).Render(context.Background(), blogIndexFile)
	if err != nil {
		log.Fatalf("failed to write blog index page: %v", err)
	}
	blogIndexFile.Close()
	fmt.Printf("Created %s at %s\n", "/words", blogIndexBuildPath)

	fmt.Printf("Generated static files to %s\n", Dirs.Build)
}

func GeneratePosts(md goldmark.Markdown, metadataContext parser.Context) []*models.Post {
	var posts = make([]*models.Post, 0)

	err := filepath.WalkDir(Dirs.Posts, func(path string, d fs.DirEntry, err error) error {
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
		postOutputPath := filepath.Join(dir, "index.html")
		postOutputFile, err := os.Create(postOutputPath)
		if err != nil {
			log.Fatalf("failed to create output post file: %v", err)
		}

		// Use templ to render the template containing the raw HTML.
		err = views.PostRoute(post).Render(context.Background(), postOutputFile)
		if err != nil {
			log.Fatalf("failed to write output file: %v", err)
		}
		postOutputFile.Close()

		fmt.Printf("Created %s at %s\n", post.Slug, postOutputPath)
		return nil
	})
	if err != nil {
		log.Fatal("There was an issue generating posts in the `filepath.WalkDir` function")
	}

	return posts
}
