package response

import "strconv"

type StatusError int

func (s StatusError) Error() string {
	return "response: fetch failed with status code " + strconv.Itoa(int(s))
}
