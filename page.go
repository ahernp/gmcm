package main

import (
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Page struct {
	Slug    string
	Content []byte
}

func markdownToHTML(args ...interface{}) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.TOC
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	s := markdown.ToHTML([]byte(fmt.Sprintf("%s", args...)), parser, renderer)

	return template.HTML(s)
}

func (p *Page) save() error {
	filename := "data/pages/" + p.Slug + ".md"
	return ioutil.WriteFile(filename, p.Content, 0600)
}

func loadPage(slug string) (*Page, error) {
	filename := "data/pages/" + slug + ".md"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Slug: slug, Content: content}, nil
}
