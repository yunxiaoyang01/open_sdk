package status

import (
	"fmt"
	"net/http"
)

type Status struct {
	Code    int
	Message string
}

func (s *Status) Error() string {
	return s.Message
}

func Error(code int, message string) error {
	return &Status{Code: code, Message: message}
}

func Errorf(code int, format string, a ...interface{}) error {
	return Error(code, fmt.Sprintf(format, a...))
}

var GetCode = func(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err := err.(type) {
	case *Status:
		return err.Code
	default:
		return http.StatusInternalServerError
	}
}

// func Code(err error) int {
// 	if err == nil {
// 		return http.StatusOK
// 	}
// 	switch err := err.(type) {
// 	case *Status:
// 		return err.Code
// 	case interface{ GRPCStatus() *status.Status }:
// 		switch c := err.GRPCStatus().Code(); c {
// 		case codes.OK:
// 			return http.StatusOK
// 		case codes.Canceled:
// 			return http.StatusRequestTimeout
// 		case codes.Unknown:
// 			return http.StatusInternalServerError
// 		case codes.InvalidArgument:
// 			return http.StatusBadRequest
// 		case codes.DeadlineExceeded:
// 			return http.StatusGatewayTimeout
// 		case codes.NotFound:
// 			return http.StatusNotFound
// 		case codes.AlreadyExists:
// 			return http.StatusConflict
// 		case codes.PermissionDenied:
// 			return http.StatusForbidden
// 		case codes.Unauthenticated:
// 			return http.StatusUnauthorized
// 		case codes.ResourceExhausted:
// 			return http.StatusTooManyRequests
// 		case codes.FailedPrecondition:
// 			return http.StatusPreconditionFailed
// 		case codes.Aborted:
// 			return http.StatusConflict
// 		case codes.OutOfRange:
// 			return http.StatusBadRequest
// 		case codes.Unimplemented:
// 			return http.StatusNotImplemented
// 		case codes.Internal:
// 			return http.StatusInternalServerError
// 		case codes.Unavailable:
// 			return http.StatusServiceUnavailable
// 		case codes.DataLoss:
// 			return http.StatusInternalServerError
// 		default:
// 			return int(c)
// 		}
// 	default:
// 		return http.StatusInternalServerError
// 	}
// }
