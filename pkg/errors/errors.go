package errors

import "encoding/json"

// swagger:model serviceError
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	b, _ := json.Marshal(e)

	return string(b)
}
