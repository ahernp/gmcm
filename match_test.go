package main

import "testing"

func TestMatch(test *testing.T) {
	matchResult := match(defaultMatchData.Input, defaultMatchData.Keys, defaultMatchData.Exclude)
	if matchResult.Output != defaultMatchData.Output {
		test.Errorf("Compare error.\nGot:\n%s\nExpected:\n%s",
			matchResult.Output,
			defaultMatchData.Output)
	}
}
