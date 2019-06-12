package main

import (
    "errors"
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/parser"
    "net/http"
    "regexp"
)

type Page struct {
    Title string
    Body  []byte
}

func markdownToHTML(args ...interface{}) template.HTML {
    extensions := parser.CommonExtensions | parser.AutoHeadingIDs
    parser := parser.NewWithExtensions(extensions)

    s := markdown.ToHTML([]byte(fmt.Sprintf("%s", args...)), parser, nil)
    return template.HTML(s)
}

var templates, err = template.New("").Funcs(template.FuncMap{"markdownToHTML": markdownToHTML}).ParseFiles("templates/edit.html", "templates/view.html")
var validPath = regexp.MustCompile("^/(edit|save|pages)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression.
}

func (p *Page) save() error {
    filename := "data/pages/" + p.Title + ".md"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := "data/pages/" + title + ".md"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    extensions := parser.CommonExtensions | parser.AutoHeadingIDs
    parser := parser.NewWithExtensions(extensions)
    html := markdown.ToHTML(p.Body, parser, nil)
    p.Body = []byte(template.HTML(html))
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/pages/"+title, http.StatusFound)
}

func redirectToHomeHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/pages/home", http.StatusFound)
}

func main() {
    http.HandleFunc("/", redirectToHomeHandler)
    http.HandleFunc("/pages/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
