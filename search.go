package main

import (
	"net/http"
	"os/exec"
	"strings"
	"text/template"
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

func highlightSubString(mainString string, subString string) string {
	mainStringLower := strings.ToLower(mainString)
	subStringLower := strings.ToLower(subString)
	subStringStart := strings.Index(mainStringLower, subStringLower)
	if subStringStart > -1 {
		subStringEnd := subStringStart + len(subString)
		return mainString[:subStringStart] + "<b>" + mainString[subStringStart:subStringEnd] + "</b>" + mainString[subStringEnd:]
	}
	return mainString
}

func getFilePathContentFromGrep(grepString string) (string, string) {
	splitString := strings.SplitAfterN(grepString, ":", 2)
	if len(splitString) > 1 {
		filePath := splitString[0][0 : len(splitString[0])-1]
		content := splitString[1]
		return filePath, content
	}
	return "", ""
}

func search(searchTerm string) {

	var nameMatches []string
	for mapPos := 0; mapPos < len(sitemap); mapPos++ {
		if strings.Contains(sitemap[mapPos].Name(), strings.ToLower(searchTerm)) {
			nameMatches = append(nameMatches, highlightSubString(sitemap[mapPos].Name(), searchTerm))
		}
	}

	grepString := "grep -i " + searchTerm + " " + pagesPath + "*"
	grepCmd := exec.Command("/bin/sh", "-c", grepString)
	grepResult, _ := grepCmd.Output()

	grepResults := strings.Split(string(grepResult[:]), "\n")
	var contentMatches []ContentMatch
	for grepPos := 0; grepPos < len(grepResults); grepPos++ {
		filePath, content := getFilePathContentFromGrep(grepResults[grepPos])
		if filePath != "" && content != "" {
			contentMatches = append(contentMatches,
				ContentMatch{
					Slug:    strings.ReplaceAll(filePath, pagesPath, ""),
					Content: highlightSubString(content, searchTerm)})
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
