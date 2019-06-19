package main

import "testing"

func TestDeduplicate(t *testing.T) {
	deduplicateResult := deduplicate(defaultDeduplicateData.Input)
	if deduplicateResult.Output != defaultDeduplicateData.Output {
		t.Errorf("Compare error.\nGot:\n%s\nExpected:\n%s", deduplicateResult.Output, defaultDeduplicateData.Output)
	}
}
