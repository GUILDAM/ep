package controller

import (
	"net/http"
	"regexp"
	"html/template"
)


var source =  "home/bas/go_home/src/app_d0659813-0363-4187-b63e-9ddd4655b3da" //os.Getenv("GOPATH") + "/src/github.com/guildam/ep/"
var templates = template.Must(template.ParseFiles(source+"home.html", source+"template/menu.html",
	source+"template/edit.html", source+"template/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$|^/$")

//MAIN
func ExecuteController() {

	http.HandleFunc("/", makeHandler(homeHandler))
	//http.HandleFunc("/import/", makeHandler(importHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	//http.Handle("/import/",  http.StripPrefix("/import/", http.FileServer(http.Dir(source + "import/"))))
	http.Handle("/scripts/",  http.StripPrefix("/scripts/", http.FileServer(http.Dir(source + "scripts/"))))
	http.Handle("/css/",  http.StripPrefix("/css/", http.FileServer(http.Dir(source + "css/"))))
	http.ListenAndServe(":8080", nil)
}
