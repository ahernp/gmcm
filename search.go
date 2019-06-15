package main

import "strings"

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
