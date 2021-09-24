package main

import (
	"embed"
	"io/fs"
)

//go:embed site/*
var siteFS embed.FS

var stripSiteFS fs.FS

func init() {
	stripSiteFS, _ = fs.Sub(siteFS, "site")
}
