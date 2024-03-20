package models

// Represents all the input/output directories for the site
type Directories struct {
	// the output directory for this app
	Build string
	// the input directory of markdown files to be used to generate `Page`s
	Content string
	// the input directory of markdown files to be used to generate blog `Post`s
	Posts string
	// the input directory of all the Templ components
	Views string
	// the input directory of all the Text templates used to generate feeds
	FeedTemplates string
}

// Our globally instantiated object of input/output directories for the site
var Dirs Directories
