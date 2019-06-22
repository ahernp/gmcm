package main

import (
	"testing"
)

const mainString = `In a hole in the ground there lived a hobbit
Not a nasty damp smelly hole filled with the ends of worms and things
nor a dry sandy hole with nothing in it to sit down on or eat
it was a hobbit-hole and that means comfort.`

var subStrings = []string{"worms", "hobbit", "in", "Comfort.", "not"}
var allMatches = [][]int{{98, 103}, {38, 44}, {0, 2}, {213, 221}, {45, 48}}

var expectedStrings = []string{
	"Not a nasty damp smelly hole filled with the ends of <b>worms</b> and things",
	"In a hole in the ground there lived a <b>hobbit</b>",
	"<b>In</b> a hole in the ground there lived a hobbit",
	"it was a hobbit-hole and that means <b>comfort.</b>",
	"<b>Not</b> a nasty damp smelly hole filled with the ends of worms and things"}

func TestHighlightSubString(test *testing.T) {
	for testPos := 0; testPos < len(allMatches); testPos++ {
		matches := [][]int{allMatches[testPos]}
		test.Run(subStrings[testPos], testHighlightSubStringFunc(
			mainString, matches, expectedStrings[testPos]))
	}
}

func testHighlightSubStringFunc(mainString string, matches [][]int, expected string) func(*testing.T) {
	return func(test *testing.T) {
		actual := highlightSubString(mainString, matches)
		if actual != expected {
			test.Errorf("Expected highlightSubString to return:\n%+v\nWhen highlighting: '%+v'\nBut got:\n%+v",
				expected,
				matches,
				actual)
		}
	}
}
