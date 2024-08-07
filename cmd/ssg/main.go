package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"bensmith.sh"
	"bensmith.sh/components"
	"bensmith.sh/routes"

	"github.com/a-h/templ"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/charmbracelet/log"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/anchor"
)

func main() {
	// True if we are in development mode
	var devMode bool
	flag.BoolVar(&devMode, "dev", false, "True if we are in development mode")
	flag.Parse()

	// create build directory for output
	if err := os.MkdirAll(bs.Dirs.Build, os.ModePerm); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// setup file writer for chroma CSS generation, we want to write it to a
	// static CSS file that targets the "chroma" class instead of using inline
	// HTML styles
	chromaCSSPath := filepath.Join(bs.Dirs.Build, "chroma.css")
	chromaCSSFile, err := os.Create(chromaCSSPath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	// make sure to wrap the styles generated by `chroma` with our custom CSS
	// layer
	chromaCSSFile.WriteString("@layer chroma {")
	defer func() {
		chromaCSSFile.WriteString("}")
		chromaCSSFile.Close()
	}()

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
				Attributer: anchor.Attributes{
					"data-anchor": "true",
				},
				// insert anchor before heading text, inside of the `<h.>` tag
				Position: anchor.Before,
				Texter:   anchor.Text("#"),
			},
			highlighting.NewHighlighting(
				highlighting.WithStyle("catppuccin-mocha"),
				// consume the io.Write we just created
				highlighting.WithCSSWriter(chromaCSSFile),
				highlighting.WithFormatOptions(
					chromahtml.TabWidth(2),
					// make sure we use class generation over inline styles
					chromahtml.WithClasses(true),
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
	feed := bs.NewFeed(posts)

	// create a single source of truth for all of our routes
	routesMap := bs.RoutesMap{
		"/": {
			Text:          "Ben Smith",
			Handler:       func() templ.Component { return routes.IndexRoute() },
			IsMainNavLink: true,
		},
		"/404/": {
			Text:    "Error Not Found",
			Handler: func() templ.Component { return routes.ErrorNotFound() },
		},
		"/pflags": {
			Text:    "pflags",
			Handler: func() templ.Component { return routes.Pflags() },
		},
		/* "/tags/": {
			Text:    "Tags",
			Handler: func() templ.Component { return routes.TagsRoute(tags) },
		}, */
		/* "/projects/": {
			Text:          "Projects",
			Handler:       func() templ.Component { return routes.ProjectsRoute() },
			IsMainNavLink: true,
		}, */
		"/words/": {
			Text:          "Writing",
			Handler:       func() templ.Component { return routes.BlogRoute(posts, tags) },
			IsMainNavLink: true,
		},
		"/words/feed/": {
			Text:          "RSS",
			Handler:       func() templ.Component { return routes.BlogFeedsRoute(posts) },
			IsMainNavLink: true,
		},
		"/words/feed.rss.xml": {
			Text:      "RSS",
			FeedType:  "rss",
			Generator: func() time.Time { return feed.Generate("/words/feed.rss.xml", "rss") },
		},
		"/words/feed.atom.xml": {
			Text:      "Atom",
			FeedType:  "atom",
			Generator: func() time.Time { return feed.Generate("/words/feed.atom.xml", "atom") },
		},
		"/words/feed.json": {
			Text:      "JSON",
			FeedType:  "json",
			Generator: func() time.Time { return feed.Generate("/words/feed.json", "json") },
		},
	}
	for _, post := range posts {
		routesMap[post.Slug] = bs.RouteEntry{
			Text:    post.Title,
			Handler: func() templ.Component { return routes.PostRoute(post) },
		}
	}
	for _, tag := range tags {
		routesMap[fmt.Sprintf("/tags/%s/", tag)] = bs.RouteEntry{
			Text:    fmt.Sprintf(`Tagged "%s"`, tag),
			Handler: func() templ.Component { return routes.TagRoute(tag, posts) },
		}
	}
	for _, page := range pages {
		isUsesPage := page.Title == "Uses"
		routesMap[page.Slug] = bs.RouteEntry{
			Text:          page.Title,
			Handler:       func() templ.Component { return routes.MarkdownPageRoute(page) },
			IsMainNavLink: isUsesPage,
		}
	}

	// initialize our sitemap data
	sitemap := bs.NewSitemap()
	// loop through all our routes and write the file to disk,
	// saving the timestamp of that route's output file generation
	// into our `sitemap.Routes`
	for route, routeEntry := range routesMap {
		var modified time.Time
		if routeEntry.Handler != nil {
			modified = generateOutputFile(route, routeEntry.Handler(), routesMap)
		} else if routeEntry.Generator != nil {
			modified = routeEntry.Generator()
		} else {
			continue
		}
		sitemap.Entries = append(sitemap.Entries, bs.SitemapEntry{
			Url:          route,
			LastModified: modified,
		})
	}
	// and write our sitemap to disk
	sitemap.Generate("/sitemap.xml")

	// Log successful completion of all the generation and exit
	log.Infof("Generated static files to `%s`", bs.Dirs.Build)
}

// Take a slug and Templ component/route to render to a static HTML file
// in our build output directory. Return the timestamp of that newly
// outputted HTML file's creation/modification for use in our Sitemap
func generateOutputFile(
	slug string,
	component templ.Component,
	routesMap bs.RoutesMap,
) time.Time {
	var dir, htmlFilePath string

	// if our `slug` contains "404", we should render a "/404.html" file instead of a directory
	// called "404/" with an "index.html" file in it. Otherwise, create a directory based upon
	// the provided `slug` and generate an `index.html` file for the route in it
	if strings.Contains(slug, "404") {
		dir = filepath.Join(bs.Dirs.Build)
		htmlFilePath = filepath.Join(dir, "404.html")
	} else {
		dir = filepath.Join(bs.Dirs.Build, slug)
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
	err = components.HTML(component).Render(ctx, file)
	if err != nil {
		log.Fatalf("failed to write blog index page: %v", err)
	}

	// log successful creation and return modified at time of output file
	log.Infof("Created %s", slug)
	info, _ := file.Stat()
	return info.ModTime()
}

// Walk the `bs.Dirs.Posts` directory and parse every markdown file found into
// a `Post` struct which we will use as content for our site. Returns
// a list of all the posts we found and a list of all the tags we found in
// those markdown posts
func generatePostsAndTags(
	md goldmark.Markdown,
	metadataContext parser.Context,
	devMode bool,
) ([]*bs.Post, []string) {
	var posts = make([]*bs.Post, 0)

	if err := filepath.WalkDir(bs.Dirs.Posts, func(path string, d fs.DirEntry, err error) error {
		// handle errors with reading the directory
		if err != nil {
			return err
		}
		// skip processing directories and files that aren't markdown
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		post, err := bs.NewPost(md, metadataContext, path)
		if err != nil {
			log.Fatalf("Failed to create new post from '%s'", path)
		}

		// if we're in dev mode show all posts. hide draft post in production
		if devMode || (!devMode && !post.Draft) {
			posts = append(posts, post)
		} else {
			log.Infof("Skipping draft Post %s", post.Slug)
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

// Walk the `bs.Dirs.Content` directory and parse every markdown file found into
// a `Page` struct which we will use as content for various pages on our site.
// Returns a list of the generated `Page`s
func generatePages(md goldmark.Markdown, metadataContext parser.Context) []*bs.Page {
	var pages = make([]*bs.Page, 0)

	if err := filepath.WalkDir(bs.Dirs.Content, func(path string, d fs.DirEntry, err error) error {
		// handle errors with reading the directory
		if err != nil {
			return err
		}
		// skip processing directories and files that aren't markdown, and don't
		// process any directories that are in the `bs.Dirs.Posts` directory
		if d.IsDir() || filepath.Ext(path) != ".md" || strings.Contains(path, bs.Dirs.Posts) {
			return nil
		}

		page, err := bs.NewPage(md, metadataContext, path)
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
