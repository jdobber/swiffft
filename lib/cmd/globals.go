package cmd

import (
	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"
	"github.com/jdobber/swiffft/lib/caches"
	sources "github.com/jdobber/swiffft/lib/sources"
)

var (
	Config  *iiifconfig.Config
	Sources []sources.Source
	Cache   caches.Cache
	Args    CommandOptions
)
