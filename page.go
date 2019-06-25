package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// PageTemplateData template data
type PageTemplateData struct {
	Page          *Page
	GlobalContext *GlobalContext
}

// Page containing Markdown text
type Page struct {
	Name    string
	Content []byte
}

const pagesPath = "data/pages/"
const mainMenuName = "main-menu"

var defaultMainMenu = []byte("[Home](pages/Home)\n\n[Task List](pages/Task List)\n")

var validPath = regexp.MustCompile("^/(edit|save|pages)/(.+)$")

var mainMenu = getMainMenu()

var viewPageTemplate = template.Must(template.New("").
	Funcs(template.FuncMap{"markdownToHTML": markdownToHTML}).
	ParseFiles("templates/view.html", "templates/base.html"))
var editPageTemplate = template.Must(
	template.ParseFiles("templates/edit.html", "templates/base.html"))

func getMainMenu() string {
	content, err := ioutil.ReadFile(pagesPath + mainMenuName)
	if err != nil {
		content = defaultMainMenu
	}
	return markdownToHTML(content)
}

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
	filename := pagesPath + page.Name
	return ioutil.WriteFile(filename, page.Content, 0644)
}

func loadPage(name string) (*Page, error) {
	filename := pagesPath + name
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Name: name, Content: content}, nil
}

func getName(request *http.Request) string {
	regexResult := validPath.FindStringSubmatch(request.URL.Path)
	if regexResult == nil {
		return ""
	}
	return regexResult[2] // Name is in remainder part of regex result
}

func viewPageHandler(writer http.ResponseWriter, request *http.Request) {
	name := getName(request)
	page, err := loadPage(name)
	if err != nil {
		http.Redirect(writer, request, "/edit/"+name, http.StatusFound)
		return
	}
	updateHistory(name)
	updatePageCache(page)
	templateData := PageTemplateData{Page: page, GlobalContext: &globalContext}
	viewPageTemplate.ExecuteTemplate(writer, "base", templateData)
}

func editPageHandler(writer http.ResponseWriter, request *http.Request) {
	name := getName(request)
	page, err := loadPage(name)
	if err != nil {
		page = &Page{Name: name}
	}
	templateData := PageTemplateData{Page: page, GlobalContext: &globalContext}
	editPageTemplate.ExecuteTemplate(writer, "base", templateData)
}

func savePageHandler(writer http.ResponseWriter, request *http.Request) {
	name := getName(request)
	content := request.FormValue("content")
	if runtime.GOOS != "windows" {
		// Remove carriage returns
		content = strings.ReplaceAll(content, "\r", "")
	}
	page := &Page{Name: name, Content: []byte(content)}
	err := page.save()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	updatePageCache(page)
	if name == mainMenuName {
		mainMenu = getMainMenu()
	}
	http.Redirect(writer, request, "/pages/"+name, http.StatusFound)
}

func redirectToHomeHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/pages/Home", http.StatusFound)
}
