package main

import "testing"

func TestCardgen(test *testing.T) {
	cardgenResult := generateCards(defaultCardgenData.Data, defaultCardgenData.Delim, defaultCardgenData.Template)
	if cardgenResult.Output != defaultCardgenData.Output {
		test.Errorf("Cardgen error.\nGot:\n%s\nExpected:\n%s",
			cardgenResult.Output,
			defaultCardgenData.Output)
	}
}
