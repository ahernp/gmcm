package main

import "testing"

func TestCompare(t *testing.T) {
	compareResult := compare(defaultCompareData.Input1, defaultCompareData.Input2)
	if compareResult.Output != defaultCompareData.Output {
		t.Errorf("Compare error.\nGot:\n%s\nExpected:\n%s", compareResult.Output, defaultCompareData.Output)
	}
}
