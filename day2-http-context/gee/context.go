package gee

import "net/http"

type H map[string]interface{}

var Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	StatusCode int
}
