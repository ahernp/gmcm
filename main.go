package main

import (
    "errors"
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
    "strings"

    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/html"
    "github.com/gomarkdown/markdown/parser"
)

type Page struct {
    Slug string
    Body []byte
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

var templates = make(map[string]*template.Template)

var validPath = regexp.MustCompile("^/(edit|save|pages)/([-a-zA-Z0-9]+)$")

func (p *Page) save() error {
    filename := "data/pages/" + p.Slug + ".md"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(slug string) (*Page, error) {
    filename := "data/pages/" + slug + ".md"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Slug: slug, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, name string, p *Page) error {
    template, ok := templates[name]
    if !ok {
        err := errors.New("Template not found -> " + name)
        return err
    }
    return template.ExecuteTemplate(w, "base", p)
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

func viewHandler(w http.ResponseWriter, r *http.Request, slug string) {
    p, err := loadPage(slug)
    if err != nil {
        http.Redirect(w, r, "/edit/"+slug, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, slug string) {
    p, err := loadPage(slug)
    if err != nil {
        p = &Page{Slug: slug}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, slug string) {
    body := r.FormValue("body")
    bodySansCarriageReturns := strings.ReplaceAll(body, "\r", "")
    p := &Page{Slug: slug, Body: []byte(bodySansCarriageReturns)}
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

func main() {
    templates["view"] = template.Must(template.New("").
        Funcs(template.FuncMap{"markdownToHTML": markdownToHTML}).
        ParseFiles("templates/view.html", "templates/base.html"))
    templates["edit"] = template.Must(
        template.ParseFiles("templates/edit.html", "templates/base.html"))
    staticFileServer := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))
    mediaFileServer := http.FileServer(http.Dir("media"))
    http.Handle("/media/", http.StripPrefix("/media/", mediaFileServer))
    http.HandleFunc("/", redirectToHomeHandler)
    http.HandleFunc("/pages/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
