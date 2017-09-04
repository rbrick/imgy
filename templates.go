package main

import "html/template"

var (
	signupTemplate  = template.Must(template.ParseFiles("templates/signup.html"))
	indexTemplate   = template.Must(template.ParseFiles("templates/index.html"))
	historyTemplate = template.Must(template.ParseFiles("templates/history-final.html"))
)
