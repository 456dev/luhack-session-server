package main

import (
	"html/template"
	"os"
	"path/filepath"
)

var htmlTemplates map[string]*template.Template

func parseTemplates() {
	htmlTemplates = make(map[string]*template.Template)
	templates := make([]string, 0)
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".html" {
			templates = append(templates, info.Name())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, templateName := range templates {
		tmpl, err := template.ParseFiles("templates/" + templateName)
		if err != nil {
			panic(err)
		}
		htmlTemplates[templateName] = tmpl
	}
}
