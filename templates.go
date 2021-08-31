package main

import (
	"embed"
	"html/template"

	log "github.com/sirupsen/logrus"
)

//go:embed index.htm
var templatesFS embed.FS

var templates *template.Template

func initTemplates() error {
	log.Debugln("Templates initialization started.")
	defer log.Debugln("Templates initialization finished.")

	var err error
	templates, err = template.ParseFS(templatesFS, "*.htm")
	if err != nil {
		return err
	}
	templatesFS = embed.FS{} // free memory

	return nil
}
