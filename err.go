package errs

import (
	"encoding/json"

	"google.golang.org/grpc/status"
)

// Err заменяет ошибку для фронта
type Err struct {
	code    errCode `json:"code"`             // HTTP-коды, 400, 500 и т.д.
	index   int     `json:"index"`            // Шестизначный индекс ошибки, где первые две цифры - номер сервиса, остальные - номер ошибки в сервисе
	message string  `json:"message"`          // Сообщение, которое принимают фронты или другие сервисы
	params  any     `json:"params,omitempty"` // Дополнительные параметры
}

func New(code errCode, index int, message string) *Err {
	return &Err{
		code:    code,
		index:   index,
		message: message,
	}
}

func (e Err) GetCode() int {
	return int(e.code)
}
func (e Err) GetIndex() int {
	return e.index
}
func (e Err) GetMessage() string {
	return e.message
}
func (e Err) GetParams() any {
	return e.params
}

func (e Err) Error() string {
	eJSON, marshalingErr := e.MarshalJSON()
	if marshalingErr != nil {
		return marshalingErr.Error()
	}

	return string(eJSON)
}

func (e Err) String() string {
	return e.Error()
}

func (e *Err) WithParams(params any) *Err {
	return &Err{
		code:    e.code,
		index:   e.index,
		message: e.message,
		params:  params,
	}
}

type jsonDto struct {
	Code    int    `json:"code"`
	Index   int    `json:"index"`
	Message string `json:"message"`
	Params  any    `json:"params,omitempty"`
}

func (e *Err) UnmarshalJSON(data []byte) error {
	dto := jsonDto{}

	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}

	e.code = errCode(dto.Code)
	e.index = dto.Index
	e.message = dto.Message
	e.params = dto.Params

	return nil
}

func (e Err) MarshalJSON() ([]byte, error) {
	dto := jsonDto{
		Code:    int(e.code),
		Index:   e.index,
		Message: e.message,
		Params:  e.params,
	}

	return json.Marshal(dto)
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

func Is(err error, index int) bool {
	if err == nil {
		return false
	}

	customErr, ok := Parse(err)
	if !ok {
		return false
	}

	return customErr.index == index
}
