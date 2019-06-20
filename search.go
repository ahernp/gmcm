package main

import (
	"html/template"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

// SearchTemplateData template context
type SearchTemplateData struct {
	SearchResults *SearchResults
	History       *[]string
}

// SearchResults contains all the matches found
type SearchResults struct {
	SearchTerm     string
	NameMatches    []string
	ContentMatches []ContentMatch
}

// ContentMatch contains a page content match
type ContentMatch struct {
	Slug    string
	Content string
}

var searchResults SearchResults
var searchTemplate = template.Must(
	template.ParseFiles("templates/search.html", "templates/base.html"))

func search(searchTerm string) {
	var nameMatches []string
	for mapPos := 0; mapPos < len(sitemap); mapPos++ {
		if strings.Contains(sitemap[mapPos].Name(), strings.ToLower(searchTerm)) {
			nameMatches = append(nameMatches, sitemap[mapPos].Name())
		}
	}

	grepString := "grep -i " + searchTerm + " " + pagesPath + "*"
	grepCmd := exec.Command("/bin/sh", "-c", grepString)
	grepResult, _ := grepCmd.Output()

	grepResults := strings.Split(string(grepResult[:]), "\n")
	var contentMatches []ContentMatch
	caseinsensitiveMatch := regexp.MustCompile(`(?i)` + searchTerm)
	for grepPos := 0; grepPos < len(grepResults); grepPos++ {
		splitString := strings.SplitAfterN(grepResults[grepPos], ":", 2)
		if len(splitString) > 1 {
			filePath := splitString[0][0 : len(splitString[0])-1]
			content := splitString[1]
			contentMatches = append(contentMatches,
				ContentMatch{
					Slug: strings.ReplaceAll(filePath, pagesPath, ""),
					// todo: Insert bold tags in positions before and after found text
					Content: caseinsensitiveMatch.ReplaceAllString(content, "<b>"+searchTerm+"</b>")})
		}
	}

	searchResults = SearchResults{SearchTerm: searchTerm, NameMatches: nameMatches, ContentMatches: contentMatches}
}

func searchHandler(writer http.ResponseWriter, request *http.Request) {
	searchTerm := request.FormValue("search")
	search(searchTerm)
	templateData := SearchTemplateData{SearchResults: &searchResults, History: &history}
	searchTemplate.ExecuteTemplate(writer, "base", templateData)
}
