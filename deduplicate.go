package main

import (
	"html/template"
	"net/http"
	"sort"
	"strings"
)

type DeduplicateTemplateData struct {
	DeduplicateData *DeduplicateData
	History         *[]string
}

type DeduplicateData struct {
	Input  string
	Output string
}

var defaultDeduplicateData = DeduplicateData{
	Input:  "Record3\nRecord4\nRecord4\nRecord1",
	Output: "Record1\nRecord3\nRecord4"}

var deduplicateData DeduplicateData

var toolsDeduplicateTemplate = template.Must(
	template.ParseFiles("templates/deduplicate.html", "templates/base.html"))

func deduplicate(input string) DeduplicateData {
	inputSansCarriageReturns := strings.ReplaceAll(input, "\r", "")
	inputRecords := strings.Split(inputSansCarriageReturns, "\n")
	sort.Strings(inputRecords)

	output := ""
	for i := 0; i < len(inputRecords)-1; i++ {
		if inputRecords[i] != inputRecords[i+1] {
			output = output + inputRecords[i] + "\n"
		}
	}

	if len(inputRecords) == 0 || len(inputRecords) == 1 {
		output = inputRecords[0]
	} else {
		output = output + inputRecords[len(inputRecords)-1] + "\n"
	}

	return DeduplicateData{Input: inputSansCarriageReturns, Output: output}
}

func deduplicateHandler(w http.ResponseWriter, r *http.Request) {
	deduplicateData = defaultDeduplicateData

	if r.Method == "POST" {
		input := r.FormValue("input")
		deduplicateData = deduplicate(input)
	}

	templateData := DeduplicateTemplateData{DeduplicateData: &deduplicateData, History: &history}
	toolsDeduplicateTemplate.ExecuteTemplate(w, "base", templateData)
}
