package main

import (
	"errors"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

// ContentMatch documents a page content match
type ContentMatch struct {
	Slug    string
	Content string
}

// SearchResults contains all the matches found
type SearchResults struct {
	SearchTerm     string
	NameMatches    []string
	ContentMatches []ContentMatch
}

var searchResults SearchResults

func search(searchTerm string) []string {
	var nameMatches []string
	re := regexp.MustCompile(`(?i)` + searchTerm)
	for i := 0; i < len(sitemap); i++ {
		if strings.Contains(sitemap[i].Name(), strings.ToLower(searchTerm)) {
			nameMatches = append(nameMatches, sitemap[i].Name())
		}
	}

	grepString := "grep -i " + searchTerm + " " + pagesPath + "*"
	grepCmd := exec.Command("/bin/sh", "-c", grepString)
	grepResult, _ := grepCmd.Output()

	grepResults := strings.Split(string(grepResult[:]), "\n")
	var contentMatches []ContentMatch
	for i := 0; i < len(grepResults); i++ {
		s := strings.SplitAfterN(grepResults[i], ":", 2)
		if len(s) > 1 {
			filePath := s[0][0 : len(s[0])-1]
			content := s[1]
			contentMatches = append(contentMatches,
				ContentMatch{
					Slug:    strings.ReplaceAll(filePath, pagesPath, ""),
					Content: re.ReplaceAllString(content, "<b>"+searchTerm+"</b>")})
		}
	}

	searchResults = SearchResults{SearchTerm: searchTerm, NameMatches: nameMatches, ContentMatches: contentMatches}
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
