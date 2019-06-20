package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
)

type TimersTemplateData struct {
	TimersData *TimersData
	History    *[]string
}

type TimersData struct {
	Timers []TimerData `json:"timers"`
}

type TimerData struct {
	Target string `json:"target"`
	Label  string `json:"label"`
	Slug   string `json:"slug"`
}

const timersFilename = "data/timers.json"

var timersData TimersData

var timersTemplate = template.Must(
	template.ParseFiles("templates/timers.html", "templates/base.html"))

func readTimers() error {
	content, _ := ioutil.ReadFile(timersFilename)
	return json.Unmarshal([]byte(content), &timersData)
}

func timersHandler(writer http.ResponseWriter, request *http.Request) {
	err := readTimers()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	templateData := TimersTemplateData{TimersData: &timersData, History: &history}
	timersTemplate.ExecuteTemplate(writer, "base", templateData)
}
