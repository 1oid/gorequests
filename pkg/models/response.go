package models

import "net/http"

type Response struct {
	MetaResponse http.Response
	StatusCode   int
	Headers      http.Header
	Body         []byte
}
