package main

import (
	"embed"
	"html/template"

	log "github.com/sirupsen/logrus"
)

var templates *template.Template

//go:embed templates/*.htm
var templatesFS embed.FS

func initTemplates() error {
	log.Debugln("Templates initialization started.")
	defer log.Debugln("Templates initialization finished.")

	var err error
	templates, err = template.ParseFS(templatesFS, "templates/index.htm")
	if err != nil {
		return err
	}

	templatesFS = embed.FS{}

	return nil
}
