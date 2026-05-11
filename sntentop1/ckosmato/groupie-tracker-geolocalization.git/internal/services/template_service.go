package services

import (
	"log"
	"text/template"
)

var (
	IndexTemplate      *template.Template
	NoResultsTemplate  *template.Template
	NoPageTemplate     *template.Template
	BadRequestTemplate *template.Template
)

func init() {
	var err error
	IndexTemplate, err = template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Fatalf("Error parsing index template: %v", err)
	}

	NoResultsTemplate, err = template.ParseFiles("web/templates/noresults.html")
	if err != nil {
		log.Fatalf("Error parsing noresults template: %v", err)
	}
	NoPageTemplate, err = template.ParseFiles("web/templates/404.html")
	if err != nil {
		log.Fatalf("Error parsing 404 template: %v", err)
	}
	BadRequestTemplate, err = template.ParseFiles("web/templates/400.html")
	if err != nil {
		log.Fatalf("Error parsing 400 template: %v", err)
	}

}
