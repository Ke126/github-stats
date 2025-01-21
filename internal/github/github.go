package github

import (
	"encoding/json"
	"net/http"

	"github.com/Ke126/github-stats/internal/response"
)

const GITHUB_API = "https://api.github.com"

type GitHubClient struct {
	Token string
}

type User struct {
	Username string `json:"login"`
	Avatar   string `json:"avatar_url"`
	Created  string `json:"created_at"`
}

// GetUser uses the /user endpoint to retrieve the user's
// username, avatar link, and datetime of account creation.
func (g *GitHubClient) GetUser() (User, error) {
	return get[User](g.Token, "/user")
}

type Repository struct {
	Name  string `json:"full_name"`
	Stars int    `json:"stargazers_count"`
}

// GetRepos uses the /user/repos endpoint to retrieve a slice
// of all repos the user has read access to. Each repo has a name
// in the format "{owner}/repo_name" and a number of stars.
func (g *GitHubClient) GetRepos() ([]Repository, error) {
	return get[[]Repository](g.Token, "/user/repos")
}

// GetLanguages uses the /repos/{owner}/{repo}/languages endpoint
// to retrieve a map of all languages and number of bytes used in that repository.
func (g *GitHubClient) GetLanguages(repoName string) (map[string]int, error) {
	return get[map[string]int](g.Token, "/repos/"+repoName+"/languages")
}

// GetContributions uses the /graphql endpoint to retrieve the
// total number of contributions made by the user in that year.
func (g *GitHubClient) GetContributions(year int) (int, error) {
	req, err := http.NewRequest(http.MethodPost, GITHUB_API+"/graphql", makeGraphQLRequestBody(year))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+g.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if !response.Ok(res.StatusCode) {
		return 0, response.StatusError{StatusCode: res.StatusCode}
	}

	contributions, err := parseGraphQLResponse(res)
	if err != nil {
		return 0, err
	}

	return contributions, nil
}

func get[T any](token string, path string) (T, error) {
	var out T
	req, err := http.NewRequest(http.MethodGet, GITHUB_API+path, nil)
	if err != nil {
		return out, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}
	defer res.Body.Close()

	if !response.Ok(res.StatusCode) {
		return out, response.StatusError{StatusCode: res.StatusCode}
	}

	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return out, err
	}

	return out, nil
}
