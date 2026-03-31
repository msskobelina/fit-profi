package controller

import (
	"io"
	"net/http"
)

type IO interface {
	Read(request interface{}, reader io.Reader) error
	Error(err error, r *http.Request, w http.ResponseWriter)
	Fatal(err error, r *http.Request, w http.ResponseWriter)
	Result(response interface{}, w http.ResponseWriter)
}

type pathParamKey string

func PathParam(r *http.Request, name string) string {
	v, _ := r.Context().Value(pathParamKey(name)).(string)
	return v
}

func PathParamKey(name string) any {
	return pathParamKey(name)
}

// ErrorResponse is the JSON error envelope returned on 4xx/5xx.
type ErrorResponse struct {
	Error string `json:"error" example:"validation error message"`
}
