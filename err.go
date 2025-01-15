package errs

import (
	"encoding/json"

	"google.golang.org/grpc/status"
)

// Err заменяет ошибку для фронта
type Err struct {
	Code    errCode `json:"code"`             // HTTP-коды, 400, 500 и т.д.
	Index   int     `json:"index"`            // Шестизначный индекс ошибки, где первые две цифры - номер сервиса, остальные - номер ошибки в сервисе
	Message string  `json:"message"`          // Сообщение, которое принимают фронты или другие сервисы
	Params  any     `json:"params,omitempty"` // Дополнительные параметры
}

func New(code errCode, index int, message string) error {
	return &Err{
		Code:    code,
		Index:   index,
		Message: message,
	}
}

func (e *Err) Error() string {
	eJSON, marshalingErr := json.Marshal(e)
	if marshalingErr != nil {
		return marshalingErr.Error()
	}

	return string(eJSON)
}

func Parse(err error) (customErr *Err, ok bool) {
	if err == nil {
		return nil, false
	}

	errStr := ""

	grpcErr, ok := status.FromError(err)
	if ok {
		errStr = grpcErr.Message()
	} else {
		errStr = err.Error()
	}

	customErr = &Err{}

	parsingErr := json.Unmarshal([]byte(errStr), &customErr)
	if parsingErr != nil {
		return nil, false
	}

	return customErr, true
}
