package gee

import (
	"fmt"
	"net/http"
)

type handler func(http.ResponseWriter, *http.Request)
type Engine struct {
	router map[string]handler
}

func New() *Engine {
	return &Engine{router: make(map[string]handler)}
}
func (engine *Engine) addRoute(method string, path string, handle handler) {
	key := method + "-" + path
	engine.router[key] = handle
}
func (engine *Engine) POST(path string, handle handler) {
	engine.addRoute("POST", path, handle)
}
func (engine *Engine) GET(path string, handle handler) {
	engine.addRoute("GET", path, handle)
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handle, ok := engine.router[key]; !ok {
		fmt.Fprintf(w, "404 NOT FOUND")
	} else {
		handle(w, req)
	}
}
func (engine *Engine) Run(address string) error {
	return http.ListenAndServe(address, engine)
}
