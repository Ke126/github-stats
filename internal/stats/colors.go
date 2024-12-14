package stats

import (
	"net/http"

	"github.com/Ke126/github-stats/internal/response"
	"gopkg.in/yaml.v3"
)

// LanguageColors retrieves the languages.yml file from
// https://github.com/github-linguist/linguist/blob/main/lib/linguist/languages.yml,
// and parses the yaml into a map from languages to hex color strings.
func LanguageColors() (map[string]string, error) {
	res, err := http.Get("https://raw.githubusercontent.com/github-linguist/linguist/refs/heads/main/lib/linguist/languages.yml")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if !response.Ok(res.StatusCode) {
		return nil, response.StatusError{StatusCode: res.StatusCode}
	}

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
