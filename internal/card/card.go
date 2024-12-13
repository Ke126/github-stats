package card

import "html/template"

func NewTemplate() (*template.Template, error) {
	return template.New("card.svg").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"mul": func(a, b int) int { return a * b },
	}).ParseFiles("templates/card.svg")
}
