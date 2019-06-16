package main

import "testing"

func TestReadHistory(t *testing.T) {
	history = readHistory()
	if history == nil {
		t.Errorf("History was not read.")
	}
}
