package response

func Ok(statusCode int) error {
	if statusCode >= 200 && statusCode <= 299 {
		return nil
	}
	return StatusError(statusCode)
}
