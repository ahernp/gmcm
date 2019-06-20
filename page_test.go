package main

import (
	"strings"
	"testing"
)

const expectedHTML = "<h1 id=\"toc_0\">Heading</h1>"

func TestMarkdownToHTML(test *testing.T) {
	markdown := "# Heading"
	html := markdownToHTML(markdown)
	if !strings.Contains(html, expectedHTML) {
		test.Errorf("Markdown '%s' converted to HTML '%s'.", markdown, html)
	}
}
