package main

import "testing"

func TestDeduplicate(test *testing.T) {
	deduplicateResult := sortAndDeduplicate(defaultDeduplicateData.Input)
	if deduplicateResult.Output != defaultDeduplicateData.Output {
		test.Errorf("Compare error.\nGot:\n%s\nExpected:\n%s",
			deduplicateResult.Output,
			defaultDeduplicateData.Output)
	}
}
