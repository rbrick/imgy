package main

import "html/template"

var (
	signInTemplate = template.Must(template.ParseFiles("templates/signin.html"))
	indexTemplate  = template.Must(template.ParseFiles("templates/index.html"))
)
