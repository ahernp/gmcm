package main

import "testing"

func TestMatch(t *testing.T) {
	matchResult := match(defaultMatchData.Input, defaultMatchData.Keys, defaultMatchData.Exclude)
	if matchResult.Output != defaultMatchData.Output {
		t.Errorf("Compare error.\nGot:\n%s\nExpected:\n%s", matchResult.Output, defaultMatchData.Output)
	}
}
