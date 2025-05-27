package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"regexp"
	"strings"
	"user-service/internal/delivery/grpc/errors"
	serviceErrors "user-service/internal/service/errors"
)

type Validatable interface {
	ValidateAll() error
}

func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		if v, ok := req.(Validatable); ok {
			if err := v.ValidateAll(); err != nil {
				return nil, errors.ParseError(serviceErrors.ApplicationError{
					StatusCode: 400,
					Code:       "Validation",
					Errors:     extractValidationErrors(err),
				})
			}
		}

		return handler(ctx, req)
	}
}

var validationRegex = regexp.MustCompile(`invalid (\w+)\.(\w+): ([^;]+)`)

func extractValidationErrors(err error) map[string]string {
	errorsMap := make(map[string][]string)

	matches := validationRegex.FindAllStringSubmatch(err.Error(), -1)
	for _, match := range matches {
		if len(match) == 4 {
			field := match[2]
			message := match[3]
			errorsMap[field] = append(errorsMap[field], message)
		}
	}

	result := make(map[string]string)
	for field, messages := range errorsMap {
		result[field] = strings.Join(messages, "; ")
	}

	if len(result) == 0 {
		result["message"] = err.Error()
	}

	return result
}
