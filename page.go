package main

import (
	"fmt"
	htmlTemplate "html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	textTemplate "text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type PageTemplateData struct {
	Page    *Page
	History *[]string
}

// Page containing Markdown text
type Page struct {
	Slug    string
	Content []byte
}

const pagesPath = "data/pages/"

var validPath = regexp.MustCompile("^/(edit|save|pages)/([-a-zA-Z0-9]+)$")
var viewPageTemplate = textTemplate.Must(textTemplate.New("").
	Funcs(textTemplate.FuncMap{"markdownToHTML": markdownToHTML}).
	ParseFiles("templates/view.html", "templates/base.html"))
var editPageTemplate = htmlTemplate.Must(
	htmlTemplate.ParseFiles("templates/edit.html", "templates/base.html"))

func markdownToHTML(args ...interface{}) string {
	// Todo: write own markdown to html converter
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.TOC
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	s := markdown.ToHTML([]byte(fmt.Sprintf("%s", args...)), parser, renderer)

	return strings.ReplaceAll(string(s), "&amp;#", "&#") // Unescape ampersands preceeding symbol number
}

func (p *Page) save() error {
	filename := pagesPath + p.Slug
	return ioutil.WriteFile(filename, p.Content, 0600)
}

func loadPage(slug string) (*Page, error) {
	filename := pagesPath + slug
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Slug: slug, Content: content}, nil
}

func getSlug(w http.ResponseWriter, r *http.Request) string {
	regexResult := validPath.FindStringSubmatch(r.URL.Path)
	if regexResult == nil {
		return ""
	}
	return regexResult[2]
}

func viewPageHandler(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(w, r)
	p, err := loadPage(slug)
	if err != nil {
		http.Redirect(w, r, "/edit/"+slug, http.StatusFound)
		return
	}
	updateHistory(slug)
	templateData := PageTemplateData{Page: p, History: &history}
	viewPageTemplate.ExecuteTemplate(w, "base", templateData)

}

func editPageHandler(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(w, r)
	p, err := loadPage(slug)
	if err != nil {
		p = &Page{Slug: slug}
	}
	templateData := PageTemplateData{Page: p, History: &history}
	editPageTemplate.ExecuteTemplate(w, "base", templateData)
}

func savePageHandler(w http.ResponseWriter, r *http.Request) {
	slug := getSlug(w, r)
	content := r.FormValue("content")
	contentSansCarriageReturns := strings.ReplaceAll(content, "\r", "")
	p := &Page{Slug: slug, Content: []byte(contentSansCarriageReturns)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/pages/"+slug, http.StatusFound)
}

func redirectToHomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/pages/home", http.StatusFound)
}
