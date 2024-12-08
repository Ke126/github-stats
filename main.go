package main

import (
	"fmt"
	"os"

	"github.com/Ke126/github-stats/internal/card"
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

	f, err := os.Create("card.svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	template, err := card.NewTemplate()
	if err != nil {
		panic(err)
	}

	err = template.Execute(f, stats)
	if err != nil {
		panic(err)
	}
}
