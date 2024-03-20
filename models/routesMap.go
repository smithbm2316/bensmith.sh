package models

import (
	"time"

	"github.com/a-h/templ"
)

// A controller that is mapped to one unique slug on our site
type RouteEntry struct {
	// The text that should be shown when this entry ends up as a link on the site
	Text string
	// Mutually exclusive with Generator`. Represents a closure that will
	// render a page for the `slug` that this entry is mapped to
	Handler func() templ.Component
	// Mutually exclusive with `Handler`. Represents a closure that will
	// generate a feed for the `slug` that this entry is mapped to
	Generator func() time.Time
	// If this entry has a `Generator`, this represents the type of feed
	// that it will generate. One of `rss`, `atom`, or `json`
	FeedType string
	// Whether or not this slug should show up on the site's MainNav menu
	IsMainNavLink bool
}

// The single source of truth for all of the routes on our site.
// Acts as an array of Controllers in a traditional MVC application
type RoutesMap map[string]RouteEntry
