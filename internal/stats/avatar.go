package stats

import (
	"encoding/base64"
	"io"
	"net/http"
)

// Base64Avatar retrieves the base64 string encoding
// of the user's avatar on GitHub.
func Base64Avatar(avatarUrl string) (string, error) {
	res, err := http.Get(avatarUrl)
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
