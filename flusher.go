package response

import (
	"net/http"
)

type Response interface {
	JSON(w http.ResponseWriter, r *http.Request)
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
	res.JSON(f.w, f.r)
}
