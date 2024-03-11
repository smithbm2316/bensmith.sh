package models

type SiteSettings struct {
	BuildDir string
	ViewsDir string
}

var Site = SiteSettings{
	BuildDir: "site",
	ViewsDir: "internal/views",
}
