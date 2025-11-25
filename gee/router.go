package gee

import (
	"fmt"
	"strings"
)

type router struct {
	tries    map[string]*node
	handlers map[string]handlerFunc
}

func newRouter() *router {
	return &router{
		tries:    make(map[string]*node),
		handlers: make(map[string]handlerFunc, 64),
	}
}

// 解析路径，构成可用节点数组
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	answer := []string{}
	for _, s := range vs {
		if s != "" {
			answer = append(answer, s)
			if s[0] == '*' {
				break
			}
		}
	}
	return answer
}
func (r *router) addRoute(method, path string, h handlerFunc) {
	key := method + "-" + path
	_, ok := r.tries[method]
	if !ok {
		r.tries[method] = &node{}
	}
	r.tries[method].insert(path, parsePattern(path), 0)
	r.handlers[key] = h
}
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	tree, ok := r.tries[method]
	if !ok {
		return nil, nil
	}
	searchPattern := parsePattern(path)
	n := tree.search(path, searchPattern, 0)
	if n != nil {
		params := make(map[string]string)
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchPattern[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchPattern[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
func (r *router) handle(c *Context) {
	//key := c.Req.Method + "-" + c.Req.URL.Path
	//if handle, ok := r.handlers[key]; !ok {
	//	fmt.Fprintf(c.Writer, "404 NOT FOUND")
	//} else {
	//	handle(c)
	//}
	n, params := r.getRoute(c.Method, c.Path)
	if n == nil {
		fmt.Fprintf(c.Writer, "404 NOT FOUND")
		return
	}
	c.Params = params
	key := c.Method + "-" + n.pattern
	r.handlers[key](c)
}
