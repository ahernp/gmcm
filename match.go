package main

import (
	"html/template"
	"net/http"
	"strings"
)

type MatchData struct {
	Input   string
	Keys    string
	Exclude bool
	Output  string
}

var defaultMatchData = MatchData{
	Input:   "Record3\nRecord4\nRecord4\nRecord1",
	Keys:    "Record4\nRecord1",
	Exclude: false,
	Output:  "Record4\nRecord4\nRecord1"}

var matchData MatchData

var toolsMatchTemplate = template.Must(
	template.ParseFiles("templates/match.html", "templates/base.html"))

func match(input string, keys string, exclude bool) MatchData {
	inputSansCarriageReturns := strings.ReplaceAll(input, "\r", "")
	keysSansCarriageReturns := strings.ReplaceAll(keys, "\r", "")
	inputRecords := strings.Split(inputSansCarriageReturns, "\n")
	keyRecords := strings.Split(keysSansCarriageReturns, "\n")

	output := ""
	for i := 0; i < len(inputRecords); i++ {
		var matchFound = false
		for j := 0; j < len(keyRecords); j++ {
			if strings.Index(inputRecords[i], keyRecords[j]) > -1 {
				matchFound = true
				break
			}
		}
		if (!matchFound && exclude) || (matchFound && !exclude) {
			output = output + inputRecords[i] + "\n"
		}
	}

	return MatchData{Input: input, Keys: keys, Exclude: exclude, Output: output}
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	matchData = defaultMatchData

	if r.Method == "POST" {
		input := r.FormValue("input")
		keys := r.FormValue("keys")
		excludeValue := r.FormValue("exclude")
		exclude := excludeValue == "checked"
		matchData = match(input, keys, exclude)
	}

	templateData = TemplateData{MatchData: &matchData, History: &history}
	toolsMatchTemplate.ExecuteTemplate(w, "base", templateData)
}
