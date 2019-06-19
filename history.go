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
	for i := 0; i < len(history); i++ {
		if history[i] != slug {
			newHistory = append(newHistory, history[i])
		}
	}
	history = newHistory[:historySize]
	writeHistory()
}
