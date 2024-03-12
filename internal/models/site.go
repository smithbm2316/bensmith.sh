package models

type SiteSettings struct {
	BuildDir   string
	ContentDir string
	PostsDir   string
	ViewsDir   string
}

var Site = SiteSettings{
	BuildDir:   "src",
	ContentDir: "internal/content",
	PostsDir:   "internal/content/words",
	ViewsDir:   "internal/views",
}
