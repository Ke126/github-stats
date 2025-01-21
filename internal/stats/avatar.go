package stats

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/Ke126/github-stats/internal/response"
)

// base64Avatar retrieves the base64 string encoding
// of the user's avatar on GitHub.
func base64Avatar(avatarUrl string) (string, error) {
	res, err := http.Get(avatarUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if !response.Ok(res.StatusCode) {
		return "", response.StatusError{StatusCode: res.StatusCode}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	b64 := base64.StdEncoding.EncodeToString(body)
	return b64, nil
}
