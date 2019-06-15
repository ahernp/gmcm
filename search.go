package main

import (
	"errors"
	"net/http"
	"strings"
)

type ContentMatch struct {
	Slug    string
	Content string
}

type SearchResults struct {
	SearchTerm     string
	NameMatches    []string
	ContentMatches []ContentMatch
}

var searchResults SearchResults

func search(searchTerm string) []string {
	nameMatches := make([]string, 1, len(sitemap))
	for i := 0; i < len(sitemap); i++ {
		if strings.Contains(sitemap[i].Name(), searchTerm) {
			nameMatches = append(nameMatches, sitemap[i].Name())
		}
	}
	searchResults = SearchResults{SearchTerm: searchTerm, NameMatches: nameMatches}
	return nil
}

func renderSearchTemplate(w http.ResponseWriter, name string, searchTerm string) error {
	template, ok := templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	search(searchTerm)
	templateData = TemplateData{SearchResults: &searchResults, History: &history}
	return template.ExecuteTemplate(w, "base", templateData)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.FormValue("search")
	renderSearchTemplate(w, "search", searchTerm)
}
