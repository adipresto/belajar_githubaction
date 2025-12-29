package exceptions

import "net/http"

func MapToHttpStatusCode(errorMsg error) int {
	var httpStatusCode int
	switch errorMsg {
		case ErrNotFound:
			httpStatusCode = http.StatusNotFound
		case ErrConflict:
			httpStatusCode = http.StatusConflict
		case ErrorBadRequest:
			httpStatusCode = http.StatusBadRequest
		default:
			httpStatusCode = http.StatusInternalServerError
	}

	return httpStatusCode
}
