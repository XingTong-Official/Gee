package gee

import (
	"net/http"
)

type handlerFunc func(c *Context)
type Engine struct {
	*RouterGroup

	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{}
	router := newRouter()
	engine.RouterGroup = &RouterGroup{engine: engine, router: router}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.RouterGroup.router.handle(c)
}
func (engine *Engine) Run(address string) error {
	return http.ListenAndServe(address, engine)
}

type RouterGroup struct {
	*router
	prefix     string
	middleWare []handlerFunc
	engine     *Engine
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	pre := r.prefix + prefix
	newGroup := &RouterGroup{
		prefix: pre,
		engine: r.engine,
		router: r.router,
	}
	r.engine.groups = append(r.engine.groups, newGroup)
	return newGroup
}
func (r *RouterGroup) addRoute(method string, path string, handler handlerFunc) {
	engine := r.engine
	p := r.prefix + path
	engine.router.addRoute(method, p, handler)
}
func (r *RouterGroup) GET(path string, handle handlerFunc) {
	r.addRoute("GET", path, handle)
}
func (r *RouterGroup) POST(path string, handle handlerFunc) {
	r.addRoute("POST", path, handle)
}
