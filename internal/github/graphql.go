package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func MakeGraphQLRequest(year int) *bytes.Reader {
	graphql := `{ "query": "query { viewer { contributionsCollection(to: \"%d-12-31T23:59:59\") { contributionCalendar { totalContributions } } } }" }`
	body := fmt.Sprintf(graphql, year)
	return bytes.NewReader([]byte(body))
}

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

func ParseGraphQLResponse(res *http.Response) (int, error) {
	var out GraphQLResponse
	err := json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return 0, err
	}
	return out.Data.Viewer.ContributionsCollection.ContributionCalendar.TotalContributions, nil
}
