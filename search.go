package main

import (
	"net/http"
	"regexp"
	"sort"
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
	Name            string
	Content         string
	NumberOfMatches int
}

var searchResults SearchResults
var searchTemplate = template.Must(
	template.ParseFiles("templates/search.html", "templates/base.html"))

var pageCache = map[string]string{}

func cacheAllPages() {
	for mapPos := 0; mapPos < len(sitemap); mapPos++ {
		page, _ := loadPage(sitemap[mapPos].Name())
		pageCache[page.Name] = string(page.Content)
	}
}

func updatePageCache(page *Page) {
	pageCache[page.Name] = string(page.Content)
}

func highlightSubString(mainString string, matches [][]int) string {
	subStringStartPos := matches[0][0]
	subStringEndPos := matches[0][1]
	lineStartPos := strings.LastIndex(mainString[:subStringStartPos], "\n") + 1
	lineEndPos := strings.Index(mainString[subStringEndPos:], "\n")
	if lineEndPos == -1 {
		lineEndPos = len(mainString)
	} else {
		lineEndPos = lineEndPos + subStringEndPos
	}
	return mainString[lineStartPos:subStringStartPos] + "<b>" + mainString[subStringStartPos:subStringEndPos] + "</b>" + mainString[subStringEndPos:lineEndPos]
}

func search(searchTerm string) {
	var nameMatches []string
	var contentMatches []ContentMatch
	caseinsensitiveMatch, err := regexp.Compile(`(?i)` + searchTerm)
	if err != nil {
		searchResults = SearchResults{
			SearchTerm:     searchTerm,
			NameMatches:    nameMatches,
			ContentMatches: contentMatches}
		return
	}

	for mapPos := 0; mapPos < len(sitemap); mapPos++ {
		if caseinsensitiveMatch.MatchString(sitemap[mapPos].Name()) {
			nameMatches = append(nameMatches, sitemap[mapPos].Name())
		}
	}

	for name, content := range pageCache {
		matches := caseinsensitiveMatch.FindAllStringIndex(content, -1)
		if len(matches) > 0 {
			contentMatches = append(contentMatches,
				ContentMatch{
					Name:            name,
					Content:         highlightSubString(content, matches),
					NumberOfMatches: len(matches)})
		}
	}

	sort.Slice(contentMatches, func(i, j int) bool {
		if contentMatches[i].NumberOfMatches == contentMatches[j].NumberOfMatches {
			return contentMatches[i].Name < contentMatches[j].Name
		}
		return contentMatches[i].NumberOfMatches > contentMatches[j].NumberOfMatches
	})

	searchResults = SearchResults{SearchTerm: searchTerm, NameMatches: nameMatches, ContentMatches: contentMatches}
}

func searchHandler(writer http.ResponseWriter, request *http.Request) {
	searchTerm := request.FormValue("search")
	search(searchTerm)
	templateData := SearchTemplateData{SearchResults: &searchResults, History: &history}
	searchTemplate.ExecuteTemplate(writer, "base", templateData)
}
