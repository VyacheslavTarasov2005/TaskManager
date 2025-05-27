package errors

type ApplicationError struct {
	StatusCode int               `json:"status_code"`
	Code       string            `json:"code"`
	Errors     map[string]string `json:"errors"`
}

func (e ApplicationError) Error() string {
	return e.Code
}
