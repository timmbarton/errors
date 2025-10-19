package errs

var (
	ErrUnknown = &Err{
		code:    ErrCodeUnknown,
		index:   0,
		message: "UNKNOWN",
	}
)
