package main

import (
	"io/ioutil"
	"strings"
)

const historySize = 20
const historyFilename = "data/history.txt"

var history = readHistory()

func readHistory() []string {
	content, err := ioutil.ReadFile(historyFilename)
	if err != nil {
		return nil
	}
	return strings.Split(string(content), "\n")
}

func writeHistory() error {
	historyAsString := strings.Join(history, "\n")
	return ioutil.WriteFile(historyFilename, []byte(historyAsString), 0600)
}

func updateHistory(slug string) {
	newHistory := []string{slug}
	for recordPos := 0; recordPos < len(history); recordPos++ {
		if history[recordPos] != slug {
			newHistory = append(newHistory, history[recordPos])
		}
	}
	history = newHistory[:historySize]
	writeHistory()
}
