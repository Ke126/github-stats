package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

const GITHUB_API = "https://api.github.com"

type GitHubClient struct {
	Token string
}

type GitHubStats struct {
	Username string
	Avatar   string

	Stars         int
	Contributions int
	Repositories  int

	Top5 []Language
}

func (g *GitHubClient) GetStats() (GitHubStats, error) {
	userInfo, err := g.FetchUserInfo()
	if err != nil {
		return GitHubStats{}, err
	}
	fmt.Println(userInfo)

	avatar, err := g.FetchAvatar(userInfo.Avatar)
	if err != nil {
		return GitHubStats{}, err
	}
	fmt.Println(avatar)

	repos, err := g.FetchRepos()
	if err != nil {
		return GitHubStats{}, err
	}

	stars := 0
	for _, repo := range repos {
		stars += repo.Stars
	}
	fmt.Println(stars)

	langBytes, err := g.FetchLanguages(repos)
	if err != nil {
		return GitHubStats{}, err
	}
	fmt.Println(langBytes)

	contributions, err := g.FetchContributions(userInfo.Created)
	if err != nil {
		return GitHubStats{}, err
	}
	fmt.Println(contributions)

	langColors, err := g.FetchLanguageColors()
	if err != nil {
		return GitHubStats{}, err
	}
	fmt.Println(langColors)

	// manual corrections before getting top 5 languages
	delete(langBytes, "ShaderLab")        // remove ShaderLab from entries
	langBytes["JavaScript"] -= 3 * 612000 // remove 3 * 612kb of static JavaScript files

	langs := g.CalculateTop5Languages(langBytes, langColors)
	fmt.Println(langs)

	return GitHubStats{
		Username: userInfo.Username,
		Avatar:   avatar,

		Stars:         stars,
		Contributions: contributions,
		Repositories:  len(repos),

		Top5: langs,
	}, nil
}

type UserInfo struct {
	Username string `json:"login"`
	Avatar   string `json:"avatar_url"`
	Created  string `json:"created_at"`
}

// FetchUserInfo uses the /user endpoint to retrieve the user's
// username, avatar link, and datetime of account creation.
func (g *GitHubClient) FetchUserInfo() (UserInfo, error) {
	req, err := http.NewRequest(http.MethodGet, GITHUB_API+"/user", nil)
	if err != nil {
		return UserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+g.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	defer res.Body.Close()

	var out UserInfo
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return UserInfo{}, err
	}

	return out, nil
}

type Repository struct {
	Name  string `json:"full_name"`
	Stars int    `json:"stargazers_count"`
}

// FetchRepos uses the /user/repos endpoint to retrieve a slice
// of all repos the user has read access to, containing each repo's
// name in the format "user/repo_name" and the number of stars.
func (g *GitHubClient) FetchRepos() ([]Repository, error) {
	req, err := http.NewRequest(http.MethodGet, GITHUB_API+"/user/repos", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+g.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var out []Repository
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// FetchLanguages uses the /repos/{user}/{repo_name}/languages endpoint
// to aggregate all the languages used by the user in all of their repos.
func (g *GitHubClient) FetchLanguages(repos []Repository) (map[string]int, error) {
	allLangs := make(map[string]int)
	for _, repo := range repos {
		fmt.Println(repo.Name)
		req, err := http.NewRequest(http.MethodGet, GITHUB_API+"/repos/"+repo.Name+"/languages", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+g.Token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		var langs map[string]int
		err = json.NewDecoder(res.Body).Decode(&langs)
		if err != nil {
			return nil, err
		}
		fmt.Println(repo, langs)

		for lang, bytes := range langs {
			allLangs[lang] += bytes
		}
	}
	return allLangs, nil
}

// struct to parse graphql json response
type GraphQLResponse struct {
	Data struct {
		Viewer struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int `json:"totalContributions"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"viewer"`
	} `json:"data"`
}

// FetchContributions uses the /graphql endpoint to retrieve the
// total number of contributions made by the user, from the datetime
// specified in start until the end of the current year.
func (g *GitHubClient) FetchContributions(start string) (int, error) {
	created, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return 0, err
	}

	contributions := 0
	// iterate over all years between the year of creation
	// and the current year inclusive
	for year := created.Year(); year <= time.Now().Year(); year++ {
		graphql := `{ "query": "query { viewer { contributionsCollection(to: \"%d-12-31T23:59:59\") { contributionCalendar { totalContributions } } } }" }`
		body := fmt.Sprintf(graphql, year)
		// fmt.Println(body)

		req, err := http.NewRequest(http.MethodPost, GITHUB_API+"/graphql", bytes.NewReader([]byte(body)))
		if err != nil {
			return 0, err
		}
		req.Header.Set("Authorization", "Bearer "+g.Token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return 0, err
		}
		defer res.Body.Close()

		var out GraphQLResponse
		err = json.NewDecoder(res.Body).Decode(&out)
		if err != nil {
			return 0, err
		}

		contributions += out.Data.Viewer.ContributionsCollection.ContributionCalendar.TotalContributions
	}

	return contributions, nil
}

// FetchAvatar retrieves the png of the user's avatar on GitHub
// and converts it to a base64 string.
// No token required.
func (g *GitHubClient) FetchAvatar(avatar string) (string, error) {
	res, err := http.Get(avatar)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	b64 := base64.StdEncoding.EncodeToString(body)
	return b64, nil
}

// FetchLanguageColors retrieves the languages.yml file from
// https://github.com/github-linguist/linguist and parses the yaml
// into a map mapping languages to hex color strings.
// No token required.
func (g *GitHubClient) FetchLanguageColors() (map[string]string, error) {
	res, err := http.Get("https://raw.githubusercontent.com/github-linguist/linguist/refs/heads/main/lib/linguist/languages.yml")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var temp map[string]struct {
		Color string `yaml:"color"`
	}
	err = yaml.NewDecoder(res.Body).Decode(&temp)
	if err != nil {
		return nil, err
	}

	// transform from map[string]struct{...} to map[string]string
	colors := make(map[string]string)
	for k, v := range temp {
		colors[k] = v.Color
	}

	return colors, nil
}

type Language struct {
	Language string
	Percent  string
	Color    string
}

func (g *GitHubClient) CalculateTop5Languages(langBytes map[string]int, langColors map[string]string) []Language {
	totalBytes := 0
	i := 0

	type MapElem struct {
		k string
		v int
	}
	s := make([]MapElem, len(langBytes))
	for k, v := range langBytes {
		totalBytes += v
		s[i] = MapElem{
			k: k,
			v: v,
		}
		i++
	}

	// sort in descending order
	sort.Slice(s, func(i, j int) bool {
		return s[i].v > s[j].v
	})

	// pick up to the top 5 languages, calculate their percentages, get their colors
	// and add to a slice
	out := make([]Language, 0, 5)
	for i := 0; i < len(s) && i < 5; i++ {
		lang := Language{
			Language: s[i].k,
			Percent:  fmt.Sprintf("%.1f%%", 100*float64(s[i].v)/float64(totalBytes)),
			Color:    langColors[s[i].k],
		}
		out = append(out, lang)
	}

	return out
}
