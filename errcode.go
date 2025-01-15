package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errCode int

const (
	ErrCodeBadRequest     = errCode(400)
	ErrCodeUnauthorized   = errCode(401)
	ErrCodeForbidden      = errCode(403)
	ErrCodeNotFound       = errCode(404)
	ErrCodeNotAllowed     = errCode(405)
	ErrCodeRequestTimeout = errCode(408)

	ErrCodeInternal       = errCode(500)
	ErrCodeNotImplemented = errCode(501)
	ErrCodeBadGateway     = errCode(502)
	ErrCodeUnknown        = errCode(520)
)

var grpcCodesMap = map[errCode]codes.Code{
	ErrCodeBadRequest:     codes.InvalidArgument,
	ErrCodeUnauthorized:   codes.Unauthenticated,
	ErrCodeForbidden:      codes.PermissionDenied,
	ErrCodeNotFound:       codes.NotFound,
	ErrCodeNotAllowed:     codes.NotFound,
	ErrCodeRequestTimeout: codes.Canceled,

	ErrCodeInternal:       codes.Internal,
	ErrCodeNotImplemented: codes.Unimplemented,
	ErrCodeBadGateway:     codes.Unavailable,
	ErrCodeUnknown:        codes.Unknown,
}

// ToGRPC returns *Err as *status.Status error
func ToGRPC(err *Err) error {
	if err == nil {
		return nil
	}

	return status.New(grpcCodesMap[err.Code], err.Error()).Err()
}
