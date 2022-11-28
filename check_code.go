package goresponse

func IsFailed(code int) bool {
	return clientError(code) || serverError(code)
}

func clientError(code int) bool {
	return code >= 400 && code < 500
}
func serverError(code int) bool {
	return code >= 500
}
