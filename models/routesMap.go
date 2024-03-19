package models

import (
	"time"

	"github.com/a-h/templ"
)

type RouteEntry struct {
	Text          string
	Handler       func() templ.Component
	Generator     func() time.Time
	FeedType      string
	IsMainNavLink bool
}

type RoutesMap map[string]RouteEntry
