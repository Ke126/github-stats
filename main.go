package main

import (
	"fmt"
	"os"
)

func main() {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		panic("GH_TOKEN is missing or unset")
	}

	client := &GitHubClient{token}

	stats, err := client.GetStats()
	if err != nil {
		panic(err)
	}

	fmt.Println(stats)

	svgCompiler, err := NewSVGCompiler("card.svg")
	if err != nil {
		panic(err)
	}

	err = svgCompiler.Compile(os.Stdout, stats)
	if err != nil {
		panic(err)
	}
}
