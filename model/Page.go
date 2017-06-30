package model

import "encoding/xml"

import (
	"html/template"
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
	Menu  template.HTML
	MenuSmall template.HTML
}

type MenuList struct {
	MenuName	xml.Name `xml:"MenuList"`
	Menu	[]Menu `xml:"Menu"`
}

type Menu struct {
	XmlName   xml.Name `xml:"Menu"`
	MenuType  string	`xml:"MenuType,attr"`
	ItemList  []string `xml:"ItemList>Item"`
}

//write pagina
func (p *Page) Save(source string, path string) error {
	filename := source + path + p.Title + ".html"
	return ioutil.WriteFile(filename, p.Body, 0600)
}