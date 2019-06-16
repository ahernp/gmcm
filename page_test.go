package main

import (
	"strings"
	"testing"
)

const expectedHTML = "<h1 id=\"toc_0\">Heading</h1>"

func TestMarkdownToHTML(t *testing.T) {
	md := "# Heading"
	html := markdownToHTML(md)
	if !strings.Contains(html, expectedHTML) {
		t.Errorf("Markdown '%s' converted to HTML '%s'.", md, html)
	}
}
