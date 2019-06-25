package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

// SitemapTemplateData template context
type SitemapTemplateData struct {
	Sitemap       *[]os.FileInfo
	GlobalContext *GlobalContext
}

var sitemap, _ = listPages() // Populate at startup to be available for searching

var sitemapTemplate = template.Must(
	template.ParseFiles("templates/sitemap.html", "templates/base.html"))

func listPages() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(pagesPath)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return files, err
}

func sitemapHandler(writer http.ResponseWriter, request *http.Request) {
	files, err := listPages()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	sitemap = files
	templateData := SitemapTemplateData{Sitemap: &sitemap, GlobalContext: &globalContext}
	sitemapTemplate.ExecuteTemplate(writer, "base", templateData)
}
