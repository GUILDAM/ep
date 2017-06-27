package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"encoding/xml"
)

var source = "C:/Users/brunoheidorn/go/src/github.com/guildam/ep/"
var templates = template.Must(template.ParseFiles(source+"home.html",
				source+"template/edit.html", source+"template/view.html"))
var validPath = regexp.MustCompile("^/(import|edit|save|view)/([a-zA-Z0-9]+)$|^/$")

type Page struct {
	Title string
	Body  []byte
}

type Menu struct {
	XMLName   xml.Name
	menuType  string
	itemList  []string
}

//valida caminho/pagina
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the third subexpression.
}

//salva pagina
func (p *Page) save(path string) error {
	filename := source + path + p.Title + ".html"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

//carrega pagina
func loadPage(path string, title string) (*Page, error) {
	filename := source + path + title + ".html"
	body, err := ioutil.ReadFile(filename)
	//bodyHtml := template.HTML(fmt.Sprintf( "%s", body))
	if err != nil {
		return &Page{Title: title, Body: body}, err
	}
	return &Page{Title: title, Body: body}, nil
}

//editar pagina
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage("edit/", title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

//executa ação de salvar
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	path := "view/"
	err := p.save(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

//carrega imports
func importHandler(w http.ResponseWriter, r *http.Request, title string) {
	//p, err := loadPage("import/", title)
	t, err := template.ParseFiles(source +"import/" + title + ".html")
	if err != nil {

		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	t.Execute(w, nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage("view/", title)
	if err != nil {

		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

//carrega home
func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	if title == "" {
		title = "home"
	}
	p, err := loadPage("", title)
	if err != nil {
		http.Redirect(w, r, "/home/", http.StatusFound)
		return
	}
	renderTemplate(w, "home", p)
}

//gerenciador de renderização de pagina
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//monta handler
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

//MAIN
func main() {
	http.HandleFunc("/", makeHandler(homeHandler))
	//http.HandleFunc("/import/", makeHandler(importHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.Handle("/import/",  http.StripPrefix("/import/", http.FileServer(http.Dir(source + "import/"))))
	http.Handle("/scripts/",  http.StripPrefix("/scripts/", http.FileServer(http.Dir(source + "scripts/"))))
	http.Handle("/css/",  http.StripPrefix("/css/", http.FileServer(http.Dir(source + "css/"))))
	http.ListenAndServe(":8080", nil)
}
