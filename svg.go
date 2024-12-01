package main

import (
	"html/template"
	"io"
)

type SVGCompiler struct {
	Template *template.Template
}

func NewSVGCompiler(templateFile string) (*SVGCompiler, error) {
	template, err := template.New(templateFile).ParseFiles(templateFile)
	if err != nil {
		return nil, err
	}
	return &SVGCompiler{template}, nil
}

func (s *SVGCompiler) Compile(writer io.Writer, stats GitHubStats) error {
	return s.Template.Execute(writer, stats)
}
