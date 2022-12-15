package goresponse

import (
	"encoding/json"
	"reflect"
)

type JSONCommon struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

type CommonOpt struct {
	Message    string
	Data       any
	Error      any
	StatusCode int
}

type Common struct {
	status     bool
	message    string
	data       any
	error      any
	statusCode int
}

func NewCommon(cfg *CommonOpt) *Common {
	return &Common{
		message:    cfg.Message,
		data:       cfg.Data,
		error:      cfg.Error,
		statusCode: cfg.StatusCode,
	}
}

func (c *Common) writeJSONSuccess() JSONCommon {
	if reflect.ValueOf(c.statusCode).IsZero() {
		c.statusCode = 200
	}

	if IsFailed(c.statusCode) {
		return c.writeJSONError()
	}

	c.status = true

	return c.toJSON()
}

func (c *Common) writeJSONError() JSONCommon {
	if reflect.ValueOf(c.statusCode).IsZero() {
		c.statusCode = 500
	}

	c.status = false
	c.data = nil

	return c.toJSON()
}

func (c *Common) Message() string {
	return c.message
}

func (c *Common) StatusCode() int {
	return c.statusCode
}

func (c *Common) toJSON() JSONCommon {
	errBytes, _ := json.Marshal(c.error)

	if string(errBytes) == `{}` {
		c.error = c.error.(error).Error()
	}

	return JSONCommon{
		Status:  c.status,
		Message: c.message,
		Data:    c.data,
		Error:   c.error,
	}
}

func (c *Common) JSONMarshal() ([]byte, error) {
	var jsonCommon JSONCommon

	if c.error != nil {
		jsonCommon = c.writeJSONError()
	} else {
		jsonCommon = c.writeJSONSuccess()
	}

	return json.Marshal(jsonCommon)
}
