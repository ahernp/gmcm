package main

import "testing"

const mainString = `In a hole in the ground there lived a hobbit
Not a nasty damp smelly hole filled with the ends of worms and things
nor a dry sandy hole with nothing in it to sit down on or eat
it was a hobbit-hole and that means comfort.`

var substrings = []string{"worms", "hobbit", "in", "Comfort.", "not"}

var expectedStrings = []string{
	"Not a nasty damp smelly hole filled with the ends of <b>worms</b> and things",
	"In a hole in the ground there lived a <b>hobbit</b>",
	"<b>In</b> a hole in the ground there lived a hobbit",
	"it was a hobbit-hole and that means <b>comfort.</b>",
	"<b>Not</b> a nasty damp smelly hole filled with the ends of worms and things"}

func TestHighlightSubString(test *testing.T) {
	for testPos := 0; testPos < len(substrings); testPos++ {
		substring := substrings[testPos]
		expected := expectedStrings[testPos]
		test.Run(substring, testHighlightSubStringFunc(mainString, substring, expected))
	}
}

func testHighlightSubStringFunc(mainString string, subString string, expected string) func(*testing.T) {
	return func(test *testing.T) {
		actual := highlightSubString(mainString, subString)
		if actual != expected {
			test.Errorf("Expected highlightSubString to return:\n%+v\nWhen highlighting: '%+v'\nBut got:\n%+v",
				expected,
				subString,
				actual)
		}
	}
}
