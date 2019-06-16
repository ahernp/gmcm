package main

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type UploadedFile struct {
	Dir  string
	File os.FileInfo
}

var uploadDirs = [...]string{"img", "code", "doc", "pres", "thumb"}
var uploadedFiles = getUploadedFiles()

var uploadTemplate = template.Must(
	template.ParseFiles("templates/upload.html", "templates/base.html"))

func getUploadedFiles() []UploadedFile {
	var uploadedFiles []UploadedFile
	for i := 0; i < len(uploadDirs); i++ {
		path := "media/" + uploadDirs[i]
		files, _ := ioutil.ReadDir(path)
		for j := 0; j < len(files); j++ {
			uploadedFiles = append(uploadedFiles, UploadedFile{Dir: uploadDirs[i], File: files[j]})
		}
	}
	return uploadedFiles
}

func renderUploadTemplate(w http.ResponseWriter) error {
	templateData = TemplateData{UploadedFiles: &uploadedFiles, History: &history}
	return uploadTemplate.ExecuteTemplate(w, "base", templateData)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		dir := r.FormValue("dir")
		r.ParseMultipartForm(10 << 20)
		sourceFile, handler, err := r.FormFile("newFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer sourceFile.Close()
		destFile, err := os.OpenFile("media/"+dir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer destFile.Close()
		io.Copy(destFile, sourceFile)
		uploadedFiles = getUploadedFiles()
	}
	renderUploadTemplate(w)
}