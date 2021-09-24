package main

import (
	"html/template"

	log "github.com/sirupsen/logrus"
)

var templates *template.Template

func initTemplates() error {
	log.Debugln("Templates initialization started.")
	defer log.Debugln("Templates initialization finished.")

	var err error
	templates, err = template.ParseFS(siteFS, "site/index.htm")
	if err != nil {
		return err
	}

	return nil
}
