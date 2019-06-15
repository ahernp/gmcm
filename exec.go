package main

import (
	"fmt"
	"os/exec"
)

func listPages() (string, error) {
	out, err := exec.Command("ls", "data/pages").Output()

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return string(out[:]), err
}
