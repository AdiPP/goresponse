package goresponse

import (
	"encoding/json"
	"github.com/lvqingan/gopager"
	"net/http"
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

type Pagination struct {
	Status  bool   `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
	Meta    meta   `json:"meta"`

	statusCode int
}

type PaginationFromLengthAwarePaginatorOpt struct {
	Message    string
	Paginator  *gopager.TLengthAwarePaginator
	Query      map[string][]string
	Error      any
	StatusCode int
}

func NewPaginationFromLengthAwarePaginator(opt *PaginationFromLengthAwarePaginatorOpt) *Pagination {
	paginator := opt.Paginator
	paginator.Appends(opt.Query)

	return &Pagination{
		Message:    opt.Message,
		Data:       opt.Paginator.Items,
		Error:      opt.Error,
		statusCode: opt.StatusCode,
		Meta:       newMeta(paginator.GetStringMap()),
	}
}

func (p *Pagination) validate() {
	if p.Error != nil {
		p.writeError()
		return
	}

	p.writeSuccess()
	return
}

func (p *Pagination) writeSuccess() {
	if reflect.ValueOf(p.statusCode).IsZero() {
		p.statusCode = 200
	}

	p.Status = true
}

func (p *Pagination) writeError() {
	if reflect.ValueOf(p.statusCode).IsZero() {
		p.statusCode = 500
	}

	p.Status = false
	p.Data = nil
}

func (p *Pagination) JSON(w http.ResponseWriter, _ *http.Request) {
	p.validate()

	var resBytes []byte
	var err error

	if resBytes, err = json.Marshal(p); err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(p.statusCode)

	if _, err = w.Write(resBytes); err != nil {
		return
	}
}
