package gee

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"sort"
	"strings"
)

type HandlerFunc func(c *Context)
type Engine struct {
	*RouterGroup

	groups        []*RouterGroup
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

func New() *Engine {
	engine := &Engine{}
	router := newRouter()
	engine.RouterGroup = &RouterGroup{engine: engine, router: router}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func Default() *Engine {
	engine := New()
	engine.Use(Recovery())
	return engine
}
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	var middleWares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(url, group.prefix) {
			middleWares = append(middleWares, group.middleWares...)
		}
	}
	c := newContext(w, req)
	c.middleWares = middleWares
	c.engine = engine
	engine.RouterGroup.router.handle(c)
}
func (engine *Engine) Run(address string) error {
	return http.ListenAndServe(address, engine)
}

type RouterGroup struct {
	*router
	prefix      string
	middleWares []HandlerFunc
	engine      *Engine
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	pre := r.prefix + prefix
	newGroup := &RouterGroup{
		prefix: pre,
		engine: r.engine,
		router: r.router,
	}
	r.engine.groups = append(r.engine.groups, newGroup)
	sort.Slice(r.engine.groups, func(i, j int) bool {
		return r.engine.groups[i].prefix < r.engine.groups[j].prefix
	})
	return newGroup
}
func (r *RouterGroup) Use(handler ...HandlerFunc) {
	r.middleWares = append(r.middleWares, handler...)
}
func (r *RouterGroup) addRoute(method string, path string, handler HandlerFunc) {
	engine := r.engine
	p := r.prefix + path
	engine.router.addRoute(method, p, handler)
}
func (r *RouterGroup) GET(path string, handle HandlerFunc) {
	r.addRoute("GET", path, handle)
}
func (r *RouterGroup) POST(path string, handle HandlerFunc) {
	r.addRoute("POST", path, handle)
}
func (r *RouterGroup) Static(relativePath string, root string) {
	r.GET(path.Join(relativePath, "/*filePath"), r.addStaticFiles(relativePath, http.Dir(root)))
}

func (r *RouterGroup) addStaticFiles(path string, fs http.FileSystem) HandlerFunc {
	pattern := r.prefix + path
	server := http.StripPrefix(pattern, http.FileServer(fs))
	return func(c *Context) {
		filename := c.Param("filePath")
		if _, err := fs.Open(filename); err != nil {
			fmt.Fprintf(c.Writer, "404 NOT FOUND")
			return
		}
		server.ServeHTTP(c.Writer, c.Req)
	}
}
