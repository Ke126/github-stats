package response

func Ok(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
