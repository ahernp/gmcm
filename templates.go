package main

import (
	"errors"
	"html/template"
	"net/http"
	"os"
)

type TemplateData struct {
	Page    *Page
	History *[]string
	Sitemap *[]os.FileInfo
}

var templateData TemplateData

var templates = make(map[string]*template.Template)

func renderTemplate(w http.ResponseWriter, name string, p *Page) error {
	template, ok := templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	templateData = TemplateData{Page: p, History: &history}
	return template.ExecuteTemplate(w, "base", templateData)
}

func renderSitemapTemplate(w http.ResponseWriter, name string, sitemap *[]os.FileInfo) error {
	template, ok := templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	templateData = TemplateData{Sitemap: sitemap, History: &history}
	return template.ExecuteTemplate(w, "base", templateData)
}
