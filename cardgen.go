package main

import (
	"html/template"
	"net/http"
	"strings"
)

// CardgenTemplateData template context
type CardgenTemplateData struct {
	CardgenData *CardgenData
	History     *[]string
}

// CardgenData form fields
type CardgenData struct {
	Data     string
	Delim    string
	Template string
	Output   string
}

const postMethod = "POST"

var defaultCardgenData = CardgenData{
	Data:     "#Name,#URL,#Description\nGoogle,www.google.com,Search engine.\nAmazon,www.amazon.co.uk,Bookshop.",
	Delim:    ",",
	Template: "<li>#Name<a href=\"https://#URL\" title=\"#Description\">#URL</a>#Description</li>",
	Output:   "<li>Google<a href=\"https://www.google.com\" title=\"Search engine.\">www.google.com</a>Search engine.</li>\n<li>Amazon<a href=\"https://www.amazon.co.uk\" title=\"Bookshop.\">www.amazon.co.uk</a>Bookshop.</li>\n"}

var cardgenData CardgenData

var cardgenTemplate = template.Must(
	template.ParseFiles("templates/cardgen.html", "templates/base.html"))

func generateCards(data string, delim string, template string) CardgenData {
	dataSansCarriageReturns := strings.ReplaceAll(data, "\r", "")
	templateSansCarriageReturns := strings.ReplaceAll(template, "\r", "")
	dataRecords := strings.Split(dataSansCarriageReturns, "\n")
	labels := strings.Split(dataRecords[0], delim)

	output := ""
	for recordPos := 1; recordPos < len(dataRecords); recordPos++ {
		currentCard := templateSansCarriageReturns
		currentData := strings.Split(dataRecords[recordPos], delim)
		for labelPos := 0; labelPos < len(currentData); labelPos++ {
			currentCard = strings.ReplaceAll(currentCard, labels[labelPos], currentData[labelPos])
		}
		output = output + currentCard + "\n"
	}
	return CardgenData{Data: data, Delim: delim, Template: template, Output: output}
}

func cardgenHandler(writer http.ResponseWriter, request *http.Request) {
	cardgenData = defaultCardgenData

	if request.Method == postMethod {
		data := request.FormValue("data")
		delim := request.FormValue("delim")
		template := request.FormValue("template")
		cardgenData = generateCards(data, delim, template)
	}

	templateData := CardgenTemplateData{CardgenData: &cardgenData, History: &history}
	cardgenTemplate.ExecuteTemplate(writer, "base", templateData)
}

func redirectToCardgenHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/tools/cardgen/", http.StatusFound)
}
