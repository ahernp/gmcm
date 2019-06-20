package main

import (
	"html/template"
	"net/http"
	"strings"
)

// MatchTemplateData template context
type MatchTemplateData struct {
	MatchData *MatchData
	History   *[]string
}

// MatchData form fields
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
	Output:  "Record4\nRecord4\nRecord1\n"}

var matchData MatchData

var matchTemplate = template.Must(
	template.ParseFiles("templates/match.html", "templates/base.html"))

func match(input string, keys string, exclude bool) MatchData {
	inputSansCarriageReturns := strings.ReplaceAll(input, "\r", "")
	keysSansCarriageReturns := strings.ReplaceAll(keys, "\r", "")
	inputRecords := strings.Split(inputSansCarriageReturns, "\n")
	keyRecords := strings.Split(keysSansCarriageReturns, "\n")

	output := ""
	for recordPos := 0; recordPos < len(inputRecords); recordPos++ {
		var matchFound = false
		for keyPos := 0; keyPos < len(keyRecords); keyPos++ {
			if strings.Index(inputRecords[recordPos], keyRecords[keyPos]) > -1 {
				matchFound = true
				break
			}
		}
		if (!matchFound && exclude) || (matchFound && !exclude) {
			output = output + inputRecords[recordPos] + "\n"
		}
	}

	return MatchData{Input: input, Keys: keys, Exclude: exclude, Output: output}
}

func matchHandler(writer http.ResponseWriter, request *http.Request) {
	matchData = defaultMatchData

	if request.Method == postMethod {
		input := request.FormValue("input")
		keys := request.FormValue("keys")
		excludeValue := request.FormValue("exclude")
		exclude := excludeValue == "checked"
		matchData = match(input, keys, exclude)
	}

	templateData := MatchTemplateData{MatchData: &matchData, History: &history}
	matchTemplate.ExecuteTemplate(writer, "base", templateData)
}
