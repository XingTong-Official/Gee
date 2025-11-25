package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	Path   string
	Method string
	Code   int
	Params map[string]string
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}
func (c *Context) Param(key string) string {
	return c.Params[key]
}
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}
func (c *Context) PostForm(key string) string {
	return c.Req.PostForm.Get(key)
}
func (c *Context) SetCode(code int) {
	c.Code = code
	c.Writer.WriteHeader(code)
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}
func (c *Context) JSON(code int, v interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetCode(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(v); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetCode(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (c *Context) HTML(code int, s string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetCode(code)
	c.Writer.Write([]byte(s))
}
func (c *Context) Data(code int, content string, data []byte) {
	c.SetHeader("Content-Type", content)
	c.SetCode(code)
	c.Writer.Write(data)
}
