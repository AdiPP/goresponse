package goresponse

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type Common struct {
	Status  bool   `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`

	statusCode int
}

type CommonOpt struct {
	Message    string
	Data       any
	Error      any
	StatusCode int
}

func NewCommon(cfg *CommonOpt) *Common {
	return &Common{
		Message:    cfg.Message,
		Data:       cfg.Data,
		Error:      cfg.Error,
		statusCode: cfg.StatusCode,
	}
}

func (c *Common) validate() {
	if c.Error != nil {
		c.writeError()
		return
	}

	c.writeSuccess()
	return
}

func (c *Common) writeSuccess() {
	if reflect.ValueOf(c.statusCode).IsZero() {
		c.statusCode = 200
	}

	c.Status = true
}

func (c *Common) writeError() {
	if reflect.ValueOf(c.statusCode).IsZero() {
		c.statusCode = 500
	}

	c.Status = false
	c.Data = nil
}

func (c *Common) JSON(w http.ResponseWriter, _ *http.Request) {
	c.validate()

	var resBytes []byte
	var err error

	if resBytes, err = json.Marshal(c); err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c.statusCode)

	if _, err = w.Write(resBytes); err != nil {
		return
	}
}
