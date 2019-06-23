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

func updateHistory(name string) {
	newHistory := []string{name}
	for recordPos := 0; recordPos < len(history); recordPos++ {
		if history[recordPos] != name {
			newHistory = append(newHistory, history[recordPos])
		}
	}
	if len(newHistory) > historySize {
		history = newHistory[:historySize]
	} else {
		history = newHistory
	}
	writeHistory()
}
