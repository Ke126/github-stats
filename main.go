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

	template, err := template.New("card.svg").ParseFiles("card.svg")
	if err != nil {
		panic(err)
	}

	err = template.Execute(os.Stdout, stats)
	if err != nil {
		panic(err)
	}
}
