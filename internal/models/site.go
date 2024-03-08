package models

type SiteSettings struct {
	BuildDir string
	ViewsDir string
}

var Site = SiteSettings{
	BuildDir: "build",
	ViewsDir: "internal/views",
}
