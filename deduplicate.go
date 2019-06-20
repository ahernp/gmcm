package main

import (
	"html/template"
	"net/http"
	"sort"
	"strings"
)

// DeduplicateTemplateData template context
type DeduplicateTemplateData struct {
	DeduplicateData *DeduplicateData
	History         *[]string
}

// DeduplicateData deduplicate form fields
type DeduplicateData struct {
	Input  string
	Output string
}

var defaultDeduplicateData = DeduplicateData{
	Input:  "Record3\nRecord4\nRecord4\nRecord1",
	Output: "Record1\nRecord3\nRecord4\n"}

var deduplicateData DeduplicateData

var deduplicateTemplate = template.Must(
	template.ParseFiles("templates/deduplicate.html", "templates/base.html"))

func sortAndDeduplicate(input string) DeduplicateData {
	inputSansCarriageReturns := strings.ReplaceAll(input, "\r", "")
	inputRecords := strings.Split(inputSansCarriageReturns, "\n")
	sort.Strings(inputRecords)

	output := ""
	for recordPos := 0; recordPos < len(inputRecords)-1; recordPos++ {
		if inputRecords[recordPos] != inputRecords[recordPos+1] {
			output = output + inputRecords[recordPos] + "\n"
		}
	}

	if len(inputRecords) == 0 || len(inputRecords) == 1 {
		output = inputRecords[0]
	} else {
		output = output + inputRecords[len(inputRecords)-1] + "\n"
	}

	return DeduplicateData{Input: inputSansCarriageReturns, Output: output}
}

func deduplicateHandler(writer http.ResponseWriter, request *http.Request) {
	deduplicateData = defaultDeduplicateData

	if request.Method == postMethod {
		input := request.FormValue("input")
		deduplicateData = sortAndDeduplicate(input)
	}

	templateData := DeduplicateTemplateData{DeduplicateData: &deduplicateData, History: &history}
	deduplicateTemplate.ExecuteTemplate(writer, "base", templateData)
}
