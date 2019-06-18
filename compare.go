package main

import (
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type CompareData struct {
	Input1 string
	Input2 string
	Output string
}

var defaultCompareData = CompareData{
	Input1: "Record1\nRecord3\nRecord4",
	Input2: "Record1\nRecord2\nRecord3",
	Output: "Results: 2 matches; 1 inserts; 1 deletes.\nI:Record2\nD:Record4"}

var compareData CompareData

var toolsCompareTemplate = template.Must(
	template.ParseFiles("templates/compare.html", "templates/base.html"))

func compare(input1 string, input2 string) CompareData {
	input1SansCarriageReturns := strings.ReplaceAll(input1, "\r", "")
	input2SansCarriageReturns := strings.ReplaceAll(input2, "\r", "")
	input1Records := strings.Split(input1SansCarriageReturns, "\n")
	input2Records := strings.Split(input2SansCarriageReturns, "\n")
	sort.Strings(input1Records)
	sort.Strings(input2Records)

	output := ""
	resultString := ""
	position1 := 0
	position2 := 0
	matchCount := 0
	insertCount := 0
	deleteCount := 0

	for position1 < len(input1Records) && position2 < len(input2Records) {
		if input1Records[position1] > input2Records[position2] {
			resultString = resultString + "I:" + input2Records[position2] + "\n"
			position2++
			insertCount++
		} else if input1Records[position1] < input2Records[position2] {
			resultString = resultString + "D:" + input1Records[position1] + "\n"
			position1++
			deleteCount++
		} else {
			position1++
			position2++
			matchCount++
		}
	}
	for position1 < len(input1Records) {
		resultString = resultString + "D:" + input1Records[position1] + "\n"
		position1++
		deleteCount++
	}
	for position2 < len(input2Records) {
		resultString = resultString + "I:" + input2Records[position2] + "\n"
		position2++
		insertCount++
	}
	summary := "Results: " + strconv.Itoa(matchCount) + " matches; " + strconv.Itoa(insertCount) + " inserts; " + strconv.Itoa(deleteCount) + " deletes.\n"
	output = summary + resultString

	return CompareData{Input1: input1SansCarriageReturns, Input2: input2SansCarriageReturns, Output: output}
}

func compareHandler(w http.ResponseWriter, r *http.Request) {
	compareData = defaultCompareData

	if r.Method == "POST" {
		input1 := r.FormValue("input1")
		input2 := r.FormValue("input2")
		compareData = compare(input1, input2)
	}

	templateData = TemplateData{CompareData: &compareData, History: &history}
	toolsCompareTemplate.ExecuteTemplate(w, "base", templateData)
}
