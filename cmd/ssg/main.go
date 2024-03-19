package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"bensmith.sh/models"
	"bensmith.sh/views"

	"github.com/a-h/templ"
	chroma "github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/anchor"
)

var Dirs = models.Directories{
	Build:         "src",
	Content:       "content",
	Posts:         "content/words",
	Views:         "views",
	FeedTemplates: "views/feeds",
}

func main() {
	// True if we are in development mode
	var devMode bool
	flag.BoolVar(&devMode, "dev", false, "True if we are in development mode")
	flag.Parse()

	// inject Dirs into nested packages
	models.Dirs = Dirs
	views.DevMode = devMode

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
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
			&anchor.Extender{
				// disable auto-adding of "anchor" class
				Attributer: anchor.Attributes{},
				// insert anchor before heading text, inside of the `<h.>` tag
				Position: anchor.Before,
			},
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
	)
	metadataContext := parser.NewContext()

	// generate posts and all their tags
	posts, tags := generatePostsAndTags(md, metadataContext, devMode)
	// generate all markdown pages
	pages := generatePages(md, metadataContext)
	// instantiate the Feed model with our generated posts
	feed := models.NewFeed(posts)

	// create a single source of truth for all of our routes
	routesMap := models.RoutesMap{
		"/": {
			Text:          "Ben Smith",
			Handler:       func() templ.Component { return views.IndexRoute() },
			IsMainNavLink: true,
		},
		"/404/": {
			Text:    "Error Not Found",
			Handler: func() templ.Component { return views.ErrorNotFound() },
		},
		"/tags/": {
			Text:    "Tags",
			Handler: func() templ.Component { return views.TagsRoute(tags) },
		},
		"/projects/": {
			Text:          "Projects",
			Handler:       func() templ.Component { return views.ProjectsRoute() },
			IsMainNavLink: true,
		},
		"/words/": {
			Text:          "Writing",
			Handler:       func() templ.Component { return views.BlogRoute(posts, tags) },
			IsMainNavLink: true,
		},
		"/words/feed/": {
			Text:          "RSS",
			Handler:       func() templ.Component { return views.BlogFeedsRoute(posts) },
			IsMainNavLink: true,
		},
		"/words/feed.rss.xml": {
			FeedType:  "rss",
			Generator: func() time.Time { return feed.Generate("/words/feed.rss.xml", "rss") },
		},
		"/words/feed.atom.xml": {
			FeedType:  "atom",
			Generator: func() time.Time { return feed.Generate("/words/feed.atom.xml", "atom") },
		},
		"/words/feed.json": {
			FeedType:  "json",
			Generator: func() time.Time { return feed.Generate("/words/feed.json", "json") },
		},
	}
	for _, post := range posts {
		routesMap[post.Slug] = models.RouteEntry{
			Text:    post.Title,
			Handler: func() templ.Component { return views.PostRoute(post) },
		}
	}
	for _, tag := range tags {
		routesMap[fmt.Sprintf("/tags/%s/", tag)] = models.RouteEntry{
			Text:    fmt.Sprintf(`Tagged "%s"`, tag),
			Handler: func() templ.Component { return views.TagsRoute(tags) },
		}
	}
	for _, page := range pages {
		isUsesPage := page.Title == "Uses"
		routesMap[page.Slug] = models.RouteEntry{
			Text:          page.Title,
			Handler:       func() templ.Component { return views.MarkdownPageRoute(page) },
			IsMainNavLink: isUsesPage,
		}
	}

	// initialize our sitemap data
	sitemap := models.NewSitemap()
	// loop through all our routes and write the file to disk
	for route, routeEntry := range routesMap {
		var modified time.Time
		if routeEntry.Handler != nil {
			modified = generateOutputFile(route, routesMap, routeEntry.Handler())
		} else if routeEntry.Generator != nil {
			modified = routeEntry.Generator()
		} else {
			continue
		}
		sitemap.Routes = append(sitemap.Routes, models.SitemapRoute{
			Url:          route,
			LastModified: modified,
		})
	}
	// and write our sitemap to disk
	sitemap.Generate("/sitemap.xml")

	// Log successful completion of all the generation and exit
	log.Printf("Generated static files to `%s`\n", Dirs.Build)
}

func generateOutputFile(
	slug string,
	routesMap models.RoutesMap,
	component templ.Component,
) time.Time {
	var dir, htmlFilePath string

	// if our `slug` contains "404", we should render a "/404.html" file instead of a directory
	// called "404/" with an "index.html" file in it. Otherwise, create a directory based upon
	// the provided `slug` and generate an `index.html` file for the route in it
	if strings.Contains(slug, "404") {
		dir = filepath.Join(Dirs.Build)
		htmlFilePath = filepath.Join(dir, "404.html")
	} else {
		dir = filepath.Join(Dirs.Build, slug)
		htmlFilePath = filepath.Join(dir, "index.html")
	}

	// create the necessary directory structure with `mkdir -p`
	if err := os.MkdirAll(dir, 0755); err != nil && err != os.ErrExist {
		log.Fatalf("failed to create dir %q: %v", dir, err)
	}

	// open and create the file writer
	file, err := os.Create(htmlFilePath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer file.Close()

	// render the specified template to the file writer
	ctx := context.WithValue(context.Background(), "currentRoute", slug)
	ctx = context.WithValue(ctx, "routesMap", routesMap)
	err = component.Render(ctx, file)
	if err != nil {
		log.Fatalf("failed to write blog index page: %v", err)
	}

	// log successful creation and return modified at time of output file
	log.Printf("Created %s\n", slug)
	info, _ := file.Stat()
	return info.ModTime()
}

func generatePostsAndTags(
	md goldmark.Markdown,
	metadataContext parser.Context,
	devMode bool,
) ([]*models.Post, []string) {
	var posts = make([]*models.Post, 0)

	if err := filepath.WalkDir(Dirs.Posts, func(path string, d fs.DirEntry, err error) error {
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

		// if we're in dev mode show all posts. hide draft post in production
		if devMode || (!devMode && !post.Draft) {
			posts = append(posts, post)
		} else {
			log.Printf("Skipping draft Post %s", post.Slug)
		}

		return nil
	}); err != nil {
		log.Fatal("There was an issue generating posts in the `filepath.WalkDir` function")
	}

	// sort the list of posts so that the newest posts are first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Published.After(posts[j].Published)
	})

	// create a map with empty structs as value since they don't allocate
	// any memory, so we can create a set of unique tags
	var tagsSet = make(map[string]struct{}, 0)
	for _, post := range posts {
		for _, tag := range post.Tags {
			if _, isPresent := tagsSet[tag]; !isPresent {
				tagsSet[tag] = struct{}{}
			}
		}
	}
	// then allocate a slice of strings which has a capacity of the length of our set
	tags := make([]string, 0, len(tagsSet))
	// append all our tags in our set to the string slice
	for tag := range tagsSet {
		tags = append(tags, tag)
	}
	// and sort it before returning it
	sort.Strings(tags)

	return posts, tags
}

func generatePages(md goldmark.Markdown, metadataContext parser.Context) []*models.Page {
	var pages = make([]*models.Page, 0)

	if err := filepath.WalkDir(Dirs.Content, func(path string, d fs.DirEntry, err error) error {
		// handle errors with reading the directory
		if err != nil {
			return err
		}
		// skip processing directories and files that aren't markdown, and don't process any directories that are in the Dirs.Posts directory
		if d.IsDir() || filepath.Ext(path) != ".md" || strings.Contains(path, Dirs.Posts) {
			return nil
		}

		page, err := models.NewPage(md, metadataContext, path)
		if err != nil {
			log.Fatalf("Failed to create new post from '%s'", path)
		}
		pages = append(pages, page)

		return nil
	}); err != nil {
		log.Fatal("There was an issue generating pages in the `filepath.WalkDir` function")
	}

	return pages
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
