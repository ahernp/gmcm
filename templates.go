package main

import (
	"errors"
	"html/template"
	"net/http"
)

type Context struct {
	Page    *Page
	History []string
}

var context Context

var templates = make(map[string]*template.Template)

func renderTemplate(w http.ResponseWriter, name string, p *Page) error {
	template, ok := templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	context = Context{Page: p, History: history}
	return template.ExecuteTemplate(w, "base", context)
}
