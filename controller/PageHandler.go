package controller

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"net/http"
	"fmt"
	"errors"
	"strings"
	. "github.com/guildam/ep/model"
)



//valida caminho/pagina
func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // The title is the third subexpression.
}

//Loads Menu
func loadMenu(editFlag bool, title string) (template.HTML, template.HTML, error){
	filename := source + "bd/menuList"
	r := MenuList{}
	var m Menu
	if editFlag {
		title = "edit"
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", "", errors.New("Menu xml not found")
	}
	err = xml.Unmarshal(file, &r)
	if err != nil {
		return "", "", errors.New("Menu not found or an error occurred")
	}
	for i:=0; len(r.Menu) > i;i++ {
		if r.Menu[i].MenuType == title {
			m = r.Menu[i]
			break
		}
	}
	var menu template.HTML
	var menuSmall template.HTML
	var menuConfig template.HTML
	var menuConfigSmall template.HTML
	for i:=0; len(m.ItemList) > i;i++ {
		if m.ItemList[i] != "save" && m.ItemList[i] != "edit" {
			menu += template.HTML("<a href='/view/" + m.ItemList[i] + "' class='w3-bar-item w3-button w3-padding-large w3-hide-small main-ucase'>" + strings.ToUpper(m.ItemList[i]) + "</a>")
			menuSmall += template.HTML("<a href='/view/" + m.ItemList[i] + "' class='w3-bar-item w3-button w3-padding-large main-ucase'>" + strings.ToUpper(m.ItemList[i]) + "</a>")
		} else if m.ItemList[i] == "save" {
			menuConfig += template.HTML("<input type='submit' value='"+ strings.ToUpper(m.ItemList[i]) +"' class='w3-bar-item w3-button w3-padding-large w3-hide-small main-ucase'></input>")
			menuConfigSmall +=  template.HTML("<input type='submit' value='"+ strings.ToUpper(m.ItemList[i]) +"' class='w3-bar-item w3-button w3-padding-large main-ucase'></a>")
		} else if m.ItemList[i] == "edit"{
			menuConfig += template.HTML("<a href='/edit/" + title + "' class='w3-bar-item w3-button w3-padding-large w3-hide-small main-ucase'>" + strings.ToUpper(m.ItemList[i]) + "</a>")
			menuSmall += template.HTML("<a href='/edit/" + title + "' class='w3-bar-item w3-button w3-padding-large main-ucase'>" + strings.ToUpper(m.ItemList[i]) + "</a>")
		}
	}
	menu += menuConfig
	menuSmall += menuConfigSmall
	return menu, menuSmall, nil
}

//Save new Menu Item
func saveMenu(title string) (error){

	filename := source + "bd/menuList"
	r := MenuList{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("Menu xml not found")
	}
	err = xml.Unmarshal(file, &r)
	if err != nil {
		return errors.New("Menu not found or an error occurred")
	}
	flagNotFound:= true
	for i:=0; len(r.Menu) > i;i++ {
		if r.Menu[i].MenuType == title {
			flagNotFound = false
			break
		}
	}
	if flagNotFound {
		//New menu item
		var newMenuItem = Menu{MenuType:title }

		for i:=0; len(r.Menu) > i;i++ {
			r.Menu[i].ItemList = append(r.Menu[i].ItemList, title)
			if r.Menu[i].MenuType == "home" {
				newMenuItem.ItemList = r.Menu[i].ItemList
				newMenuItem.ItemList = append(newMenuItem.ItemList, "edit")
			}
		}
		r.Menu = append(r.Menu, newMenuItem)

		output, err := xml.MarshalIndent(r, "  ", "    ")
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return err
		}
		return ioutil.WriteFile(filename, output, 0600)
	}
	return nil
}

//carrega pagina
func loadPage(path string, title string, editFlag bool) (*Page, error) {
	filename := source + path + title + ".html"
	menu, menuSmall, err := loadMenu(editFlag, title)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return &Page{Title: title, Body: body, Menu: menu, MenuSmall: menuSmall}, err
	}
	return &Page{Title: title, Body: body, Menu: menu, MenuSmall: menuSmall}, nil
}

//editar pagina
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage("view/", title, true)
	if err != nil {
		p = &Page{Title: title, Menu: p.Menu}
	}
	renderTemplate(w, "edit", p)
}

//salvar edição
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	path := "view/"
	err := p.Save(source, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = saveMenu(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

//loads views
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage("view/", title, false)
	if err != nil {

		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

//carrega imports depreciated code
/*func importHandler(w http.ResponseWriter, r *http.Request, title string) {
	//p, err := loadPage("import/", title)
	t, err := template.ParseFiles(source +"import/" + title + ".html")
	if err != nil {

		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	t.Execute(w, nil)
}*/

//carrega home
func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	if title == "" {
		title = "home"
	}
	p, err := loadPage("", title, false)
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

//build handler
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