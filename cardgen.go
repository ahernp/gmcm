package main

import (
	"html/template"
	"net/http"
	"strings"
)

type CardgenTemplateData struct {
	CardgenData *CardgenData
	History     *[]string
}

type CardgenData struct {
	Data     string
	Delim    string
	Template string
	Output   string
}

var defaultCardgenData = CardgenData{
	Data:     "#Name,#URL,#Description\nGoogle,www.google.com,Search engine.\nAmazon,www.amazon.co.uk,Bookshop.",
	Delim:    ",",
	Template: "<li>#Name<a href=\"https://#URL\" title=\"#Description\">#URL</a>#Description</li>",
	Output:   "<li>Google<a href=\"https://www.google.com\" title=\"Search engine.\">www.google.com</a>Search engine.</li>\n<li>Amazon<a href=\"https://www.amazon.co.uk\" title=\"Bookshop.\">www.amazon.co.uk</a>Bookshop.</li>"}

var cardgenData CardgenData

var toolsCardgenTemplate = template.Must(
	template.ParseFiles("templates/cardgen.html", "templates/base.html"))

func generateCards(data string, delim string, template string) CardgenData {
	dataSansCarriageReturns := strings.ReplaceAll(data, "\r", "")
	templateSansCarriageReturns := strings.ReplaceAll(template, "\r", "")
	dataRecords := strings.Split(dataSansCarriageReturns, "\n")
	labels := strings.Split(dataRecords[0], delim)

	output := ""
	for i := 1; i < len(dataRecords); i++ {
		currentCard := templateSansCarriageReturns
		currentData := strings.Split(dataRecords[i], delim)
		for j := 0; j < len(currentData); j++ {

			currentCard = strings.ReplaceAll(currentCard, labels[j], currentData[j])
		}
		output = output + currentCard + "\n"
	}
	return CardgenData{Data: data, Delim: delim, Template: template, Output: output}
}

func cardgenHandler(w http.ResponseWriter, r *http.Request) {
	cardgenData = defaultCardgenData

	if r.Method == "POST" {
		data := r.FormValue("data")
		delim := r.FormValue("delim")
		template := r.FormValue("template")
		cardgenData = generateCards(data, delim, template)
	}

	templateData := CardgenTemplateData{CardgenData: &cardgenData, History: &history}
	toolsCardgenTemplate.ExecuteTemplate(w, "base", templateData)
}

func redirectToCardgenHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tools/cardgen/", http.StatusFound)
}
