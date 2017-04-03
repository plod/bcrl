package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

//TemplateData is a struct for passing through to templates
type TemplateData struct {
	PageTitle string
	Error     bool
	OtherData interface{}
	Year      int
}

func notImplemented() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func home(w http.ResponseWriter, req *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	generateHTML(w, TemplateData{"", false, "The time is: " + tm, 1974}, "layout.tpl")
}

func generateHTML(writer http.ResponseWriter, td TemplateData, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, filepath.Join("templates", fmt.Sprintf("%s.html", file)))
	}
	templates := template.Must(template.ParseFiles(files...))
	td.Year = time.Now().Year()
	err := templates.ExecuteTemplate(writer, "layout", td)
	if err != nil {
		log.Println(yellow("WARN"), err)
	}
}
