package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ke126/github-stats/internal/card"
	"github.com/Ke126/github-stats/internal/github"
	"github.com/Ke126/github-stats/internal/stats"
)

func main() {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		log.Fatal("GH_TOKEN is missing or unset")
	}

	ghClient := &github.GitHubClient{Token: token}
	ghStats := &stats.GitHubStats{Client: ghClient}

	stats, err := ghStats.AllStats()
	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Println(stats)

	f, err := os.Create("card.svg")
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer f.Close()

	template, err := card.NewTemplate()
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = template.Execute(f, stats)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
