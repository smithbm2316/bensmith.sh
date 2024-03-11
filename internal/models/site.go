package models

type SiteSettings struct {
	BuildDir string
	ViewsDir string
}

var Site = SiteSettings{
	BuildDir: "src",
	ViewsDir: "internal/views",
}
