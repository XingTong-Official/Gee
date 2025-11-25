package gee

import (
	"net/http"
)

type handlerFunc func(c *Context)
type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}
func (engine *Engine) addRoute(method string, pattern string, handler handlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}
func (engine *Engine) POST(path string, handle handlerFunc) {
	engine.addRoute("POST", path, handle)
}
func (engine *Engine) GET(path string, handle handlerFunc) {
	engine.addRoute("GET", path, handle)
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
func (engine *Engine) Run(address string) error {
	return http.ListenAndServe(address, engine)
}
