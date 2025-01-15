package errs

var (
	ErrUnknown = &Err{
		Code:    ErrCodeUnknown,
		Index:   0,
		Message: "UNKNOWN",
	}
)
