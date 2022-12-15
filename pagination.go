package goresponse

import (
	"encoding/json"
	"github.com/lvqingan/gopager"
	"reflect"
)

type meta struct {
	Pagination map[string]interface{} `json:"pagination"`
}

func newMeta(pagination map[string]interface{}) meta {
	delete(pagination, "data")
	delete(pagination, "prev_page_url")

	return meta{
		Pagination: pagination,
	}
}

type JSONPagination struct {
	Status  bool   `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Meta    meta   `json:"meta"`
}

type PaginationFromLengthAwarePaginatorOpt struct {
	Message    string
	Paginator  *gopager.TLengthAwarePaginator
	Query      map[string][]string
	Error      any
	StatusCode int
}

type Pagination struct {
	status     bool
	message    string
	data       any
	error      any
	meta       meta
	statusCode int
}

func NewPaginationFromLengthAwarePaginator(opt *PaginationFromLengthAwarePaginatorOpt) *Pagination {
	paginator := opt.Paginator
	paginator.Appends(opt.Query)

	return &Pagination{
		message:    opt.Message,
		data:       opt.Paginator.Items,
		error:      opt.Error,
		statusCode: opt.StatusCode,
		meta:       newMeta(paginator.GetStringMap()),
	}
}

func (p *Pagination) writeJSONSuccess() JSONPagination {
	if reflect.ValueOf(p.statusCode).IsZero() {
		p.statusCode = 200
	}

	if IsFailed(p.statusCode) {
		return p.writeJSONError()
	}

	p.status = true

	return p.toJSON()
}

func (p *Pagination) writeJSONError() JSONPagination {
	if reflect.ValueOf(p.statusCode).IsZero() {
		p.statusCode = 500
	}

	p.status = false
	p.data = nil

	return p.toJSON()
}

func (p *Pagination) Message() string {
	return p.message
}

func (p *Pagination) StatusCode() int {
	return p.statusCode
}

func (p *Pagination) toJSON() JSONPagination {
	errBytes, _ := json.Marshal(p.error)

	if string(errBytes) == `{}` {
		p.error = p.error.(error).Error()
	}

	return JSONPagination{
		Status:  p.status,
		Message: p.message,
		Data:    p.data,
		Error:   p.error,
		Meta:    p.meta,
	}
}

func (p *Pagination) JSONMarshal() ([]byte, error) {
	var jsonCommon JSONPagination

	if p.error != nil {
		jsonCommon = p.writeJSONError()
	} else {
		jsonCommon = p.writeJSONSuccess()
	}

	return json.Marshal(jsonCommon)
}
