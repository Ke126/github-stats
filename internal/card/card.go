package card

import "html/template"

func NewTemplate() (*template.Template, error) {
	return template.New("card.svg").ParseFiles("templates/card.svg")
}
