package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"sort"
)

// TimersTemplateData template context
type TimersTemplateData struct {
	TimersData *TimersData
	History    *[]string
}

// TimersData template data
type TimersData struct {
	Timers []TimerData `json:"timers"`
}

// TimerData attributes of a single timer
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
	err := json.Unmarshal([]byte(content), &timersData)
	if err == nil {
		sort.Slice(timersData.Timers, func(i, j int) bool {
			return timersData.Timers[i].Target < timersData.Timers[j].Target
		})
	}
	return err
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
