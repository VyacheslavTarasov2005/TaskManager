package errors

import (
	"errors"
	"google.golang.org/grpc/codes"
	appErrors "user-service/internal/service/errors"
)

func ParseError(err error) codes.Code {
	var applicationError appErrors.ApplicationError
	switch {
	case errors.As(err, &applicationError):
		var appErr appErrors.ApplicationError
		errors.As(err, &appErr)
		switch appErr.StatusCode {
		case 400:
			return codes.InvalidArgument
		case 401:
			return codes.Unauthenticated
		case 403:
			return codes.PermissionDenied
		case 404:
			return codes.NotFound
		case 409:
			return codes.AlreadyExists
		default:
			return codes.Internal
		}
	default:
		return codes.Internal
	}
}
