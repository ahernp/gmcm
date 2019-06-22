package main

import (
	"fmt"
	htmlTemplate "html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	textTemplate "text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// PageTemplateData template data
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
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.TOC
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	htmlString := markdown.ToHTML([]byte(fmt.Sprintf("%s", args...)), parser, renderer)

	// Remove escaping of ampersands preceeding a symbol number
	return strings.ReplaceAll(string(htmlString), "&amp;#", "&#")
}

func (page *Page) save() error {
	filename := pagesPath + page.Slug
	return ioutil.WriteFile(filename, page.Content, 0644)
}

func loadPage(slug string) (*Page, error) {
	filename := pagesPath + slug
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Slug: slug, Content: content}, nil
}

func getSlug(request *http.Request) string {
	regexResult := validPath.FindStringSubmatch(request.URL.Path)
	if regexResult == nil {
		return ""
	}
	return regexResult[2] // Slug is in remainder part of regex result
}

func viewPageHandler(writer http.ResponseWriter, request *http.Request) {
	slug := getSlug(request)
	page, err := loadPage(slug)
	if err != nil {
		http.Redirect(writer, request, "/edit/"+slug, http.StatusFound)
		return
	}
	updateHistory(slug)
	updatePageCache(page)
	templateData := PageTemplateData{Page: page, History: &history}
	viewPageTemplate.ExecuteTemplate(writer, "base", templateData)
}

func editPageHandler(writer http.ResponseWriter, request *http.Request) {
	slug := getSlug(request)
	page, err := loadPage(slug)
	if err != nil {
		page = &Page{Slug: slug}
	}
	templateData := PageTemplateData{Page: page, History: &history}
	editPageTemplate.ExecuteTemplate(writer, "base", templateData)
}

func savePageHandler(writer http.ResponseWriter, request *http.Request) {
	slug := getSlug(request)
	content := request.FormValue("content")
	if runtime.GOOS != "windows" {
		// Remove carriage returns
		content = strings.ReplaceAll(content, "\r", "")
	}
	page := &Page{Slug: slug, Content: []byte(content)}
	err := page.save()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	updatePageCache(page)
	http.Redirect(writer, request, "/pages/"+slug, http.StatusFound)
}

func redirectToHomeHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/pages/home", http.StatusFound)
}
