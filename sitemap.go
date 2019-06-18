package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type SitemapTemplateData struct {
	Sitemap *[]os.FileInfo
	History *[]string
}

var sitemap, _ = listPages()

var sitemapTemplate = template.Must(
	template.ParseFiles("templates/sitemap.html", "templates/base.html"))

func listPages() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(pagesPath)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return files, err
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	files, err := listPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sitemap = files
	templateData := SitemapTemplateData{Sitemap: &sitemap, History: &history}
	sitemapTemplate.ExecuteTemplate(w, "base", templateData)
}
