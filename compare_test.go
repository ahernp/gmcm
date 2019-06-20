package main

import "testing"

func TestCompare(test *testing.T) {
	compareResult := compare(defaultCompareData.Input1, defaultCompareData.Input2)
	if compareResult.Output != defaultCompareData.Output {
		test.Errorf("Compare error.\nGot:\n%s\nExpected:\n%s",
			compareResult.Output,
			defaultCompareData.Output)
	}
}
