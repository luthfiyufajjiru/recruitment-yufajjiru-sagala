package customerror

import (
	"strings"
)

type (
	HttpError struct {
		Message    string
		StatusCode int
		Err        error // put this as detail error
	}
)

func (h *HttpError) Error() string {
	var sb strings.Builder
	sb.WriteString("message: ")
	sb.WriteString(h.Message)
	if h.Err != nil {
		sb.WriteString(h.Err.Error())
	}
	return sb.String()
}
