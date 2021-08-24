package main

import (
	"fmt"
	"io"
	"net/http"
)

type Request interface {
	Method(string) Request
	URL(string) Request
	Body(io.Reader) Request
	Build() (*http.Request, error)
}

type request struct {
	method    string
	url       string
	body      io.Reader
	lastError error
}

func NewRequest() Request {
	return request{}
}

func (r request) Method(method string) Request {
	switch method {
	case http.MethodGet:
	case http.MethodHead:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodPatch:
	case http.MethodDelete:
	case http.MethodConnect:
	case http.MethodOptions:
	case http.MethodTrace:
	default:
		r.lastError = fmt.Errorf("Invalid method %q provided", method)
	}
	r.method = method
	return r
}

func (r request) URL(url string) Request {
	if url == "" {
		r.lastError = fmt.Errorf("URL must not be empty")
	}
	r.url = url
	return r
}

func (r request) Body(body io.Reader) Request {
	r.body = body
	return r
}

func (r request) Build() (*http.Request, error) {
	if r.lastError != nil {
		return nil, r.lastError
	}
	req, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return nil, err
	}
	return req, nil
}
