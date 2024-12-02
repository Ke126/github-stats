package stats

import (
	"time"

	"github.com/Ke126/gh-stats/internal/github"
)

type GitHubStats struct {
	Client GitHubGetter
}

type GitHubGetter interface {
	GetContributions(year int) (int, error)
	GetLanguages(repoName string) (map[string]int, error)
	GetRepos() ([]github.Repository, error)
	GetUser() (github.User, error)
}

type Stats struct {
	Username string
	Avatar   string

	Stars         int
	Contributions int
	Repositories  int

	Top5 []Language
}

type Language struct {
	Language string
	Percent  string
	Color    string
}

func (s *GitHubStats) AllStats() (Stats, error) {
	userInfo, err := s.Client.GetUser()
	if err != nil {
		return Stats{}, err
	}

	avatar, err := Base64Avatar(userInfo.Avatar)
	if err != nil {
		return Stats{}, err
	}

	repos, err := s.Client.GetRepos()
	if err != nil {
		return Stats{}, err
	}

	stars := 0
	langs := make(map[string]int)
	for _, repo := range repos {
		stars += repo.Stars
		lang, err := s.Client.GetLanguages(repo.Name)
		if err != nil {
			return Stats{}, err
		}
		for k, v := range lang {
			langs[k] += v
		}
	}

	created, err := time.Parse(time.RFC3339, userInfo.Created)
	if err != nil {
		return Stats{}, err
	}

	contributions := 0
	// iterate over all years between the year of creation
	// and the current year inclusive
	for year := created.Year(); year <= time.Now().Year(); year++ {
		yearContributions, err := s.Client.GetContributions(year)
		if err != nil {
			return Stats{}, err
		}
		contributions += yearContributions
	}
	colors, err := LanguageColors()
	if err != nil {
		return Stats{}, err
	}

	// manual corrections before getting top 5 languages
	delete(langs, "ShaderLab")        // remove ShaderLab from entries
	langs["JavaScript"] -= 3 * 612000 // remove 3 * 612kb of static JavaScript files

	top5 := Top5Languages(langs, colors)

	return Stats{
		Username: userInfo.Username,
		Avatar:   avatar,

		Stars:         stars,
		Contributions: contributions,
		Repositories:  len(repos),

		Top5: top5,
	}, nil
}
