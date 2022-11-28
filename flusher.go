package goresponse

import (
	"net/http"
)

type Response interface {
	StatusCode() int
	JSONMarshal() ([]byte, error)
}

type Flusher struct {
	w http.ResponseWriter
	r *http.Request
}

func NewFlusher(writer http.ResponseWriter, req *http.Request) *Flusher {
	return &Flusher{
		w: writer,
		r: req,
	}
}

func (f *Flusher) JSON(res Response) {
	var resBytes []byte
	var err error

	if resBytes, err = res.JSONMarshal(); err != nil {
		return
	}

	f.w.Header().Set("Content-Type", "application/json")
	f.w.WriteHeader(res.StatusCode())

	if _, err = f.w.Write(resBytes); err != nil {
		return
	}

}
