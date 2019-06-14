package main

import (
	"io/ioutil"
	"strings"
)

const historySize = 20

var history = make([]string, historySize)

func readHistory() []string {
	filename := "data/history.txt"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return strings.Split(string(content), "\n")
}

func writeHistory() error {
	filename := "data/history.txt"
	historyAsString := strings.Join(history, "\n")
	return ioutil.WriteFile(filename, []byte(historyAsString), 0600)
}

func updateHistory(slug string) {
	newHistory := make([]string, 1, historySize)
	newHistory[0] = slug
	for i := 0; i < len(history); i++ {
		if history[i] != slug {
			newHistory = append(newHistory, history[i])
		}
	}
	history = newHistory[:historySize]
	writeHistory()
}
