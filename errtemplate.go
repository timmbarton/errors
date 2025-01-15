package errs

// ErrTemplate - шаблон для ошибки, для которой дополнительно можно указать параметры в месте возникновения
type ErrTemplate struct {
	Code    errCode `json:"code"`    // HTTP-коды, 400, 500 и т.д.
	Index   int     `json:"index"`   // Шестизначный индекс ошибки, где первые две цифры - номер сервиса, остальные - номер ошибки в сервисе
	Message string  `json:"message"` // Сообщение, которое принимают фронты или другие сервисы
}

func NewTemplate(code errCode, index int, message string) *ErrTemplate {
	return &ErrTemplate{
		Code:    code,
		Index:   index,
		Message: message,
	}
}

func (t *ErrTemplate) New(params any) error {
	return &Err{
		Code:    t.Code,
		Index:   t.Index,
		Message: t.Message,
		Params:  params,
	}
}
