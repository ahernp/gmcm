package main

import "testing"

func TestReadHistory(test *testing.T) {
	history = readHistory()
	if history == nil {
		test.Errorf("History was not read.")
	}
}
