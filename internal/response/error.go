package response

import "strconv"

type StatusError struct {
	StatusCode int
}

var _ error = StatusError{}

func (s StatusError) Error() string {
	return "fetch failed with status code " + strconv.Itoa(s.StatusCode)
}
