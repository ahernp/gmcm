package main

import (
	"net/http"
	"regexp"
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

var pageCache = map[string]string{}

func cacheAllPages() {
	for mapPos := 0; mapPos < len(sitemap); mapPos++ {
		page, _ := loadPage(sitemap[mapPos].Name())
		pageCache[page.Slug] = string(page.Content)
	}
}

func updatePageCache(page *Page) {
	pageCache[page.Slug] = string(page.Content)
}

func highlightSubString(mainString string, subString string) string {
	mainStringLower := strings.ToLower(mainString)
	subStringLower := strings.ToLower(subString)
	subStringStartPos := strings.Index(mainStringLower, subStringLower)
	if subStringStartPos > -1 {
		subStringEndPos := subStringStartPos + len(subString)
		lineStartPos := strings.LastIndex(mainString[:subStringStartPos], "\n") + 1
		lineEndPos := strings.Index(mainString[subStringEndPos:], "\n")
		if lineEndPos == -1 {
			lineEndPos = len(mainString)
		} else {
			lineEndPos = lineEndPos + subStringEndPos
		}
		return mainString[lineStartPos:subStringStartPos] + "<b>" + mainString[subStringStartPos:subStringEndPos] + "</b>" + mainString[subStringEndPos:lineEndPos]
	}
	return mainString
}

func search(searchTerm string) {
	caseinsensitiveMatch := regexp.MustCompile(`(?i)` + searchTerm)

	var nameMatches []string
	for mapPos := 0; mapPos < len(sitemap); mapPos++ {
		if caseinsensitiveMatch.MatchString(sitemap[mapPos].Name()) {
			nameMatches = append(nameMatches, sitemap[mapPos].Name())
		}
	}

	var contentMatches []ContentMatch
	for slug, content := range pageCache {
		if caseinsensitiveMatch.MatchString(content) {
			contentMatches = append(contentMatches,
				ContentMatch{
					Slug:    slug,
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
