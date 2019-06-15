package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func listPages() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir("data/pages")

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return files, err
}
