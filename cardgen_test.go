package main

import "testing"

func TestCardgen(t *testing.T) {
	cardgenResult := generateCards(defaultCardgenData.Data, defaultCardgenData.Delim, defaultCardgenData.Template)
	if cardgenResult.Output != defaultCardgenData.Output {
		t.Errorf("Cardgen error.\nGot:\n%s\nExpected:\n%s", cardgenResult.Output, defaultCardgenData.Output)
	}
}
