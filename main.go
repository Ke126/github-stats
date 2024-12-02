package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/Ke126/github-stats/internal/github"
	"github.com/Ke126/github-stats/internal/stats"
)

func main() {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		panic("GH_TOKEN is missing or unset")
	}

	ghClient := &github.GitHubClient{Token: token}
	ghStats := &stats.GitHubStats{Client: ghClient}

	stats, err := ghStats.AllStats()
	if err != nil {
		panic(err)
	}

	fmt.Println(stats)

	err = os.MkdirAll("_site", 0750)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("_site/card.svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	template := template.Must(template.New("template.svg").ParseFiles("template.svg"))

	err = template.Execute(f, stats)
	if err != nil {
		panic(err)
	}
}
