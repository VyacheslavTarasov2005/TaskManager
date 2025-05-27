package errors

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user-service/internal/delivery/grpc/pb"
	appErrors "user-service/internal/service/errors"
)

func ParseError(err error) error {
	var appErr appErrors.ApplicationError
	if errors.As(err, &appErr) {
		st := status.New(convertHTTPToGRPCCode(appErr.StatusCode), appErr.Error())

		detail := &pb.ErrorDetail{
			Code:   appErr.Code,
			Errors: appErr.Errors,
		}

		stWithDetails, err := st.WithDetails(detail)
		if err != nil {
			return st.Err()
		}

		return stWithDetails.Err()
	}
	return status.Error(codes.Internal, "internal server error")
}

func convertHTTPToGRPCCode(httpCode int) codes.Code {
	switch httpCode {
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
}
