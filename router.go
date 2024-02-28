package mule

import "strings"

type IRouter interface {
	// Add
	Handle(string, string, ...HandleFunc) IRouter
	Get(string, ...HandleFunc) IRouter
	Post(string, ...HandleFunc) IRouter

	// Query
	MatchHandlersChain(string, string) (HandlersChain, bool)
}

type Router struct {
	pathHandlerMap map[string]map[string]HandlersChain // map[path][method]HandlersChain
}

func NewRouter() IRouter {
	return &Router{pathHandlerMap: make(map[string]map[string]HandlersChain)}
}

func (r *Router) Handle(s string, s2 string, handleFunc ...HandleFunc) IRouter {
	//TODO implement me
	panic("implement me")
}

func (r *Router) initPathMap(p string) map[string]HandlersChain {
	pathMap, pathOK := r.pathHandlerMap[p]
	if !pathOK {
		pathMap = make(map[string]HandlersChain)
		r.pathHandlerMap[p] = pathMap
	}
	return pathMap
}

func (r *Router) Get(p string, handleFunc ...HandleFunc) IRouter {
	pathMap := r.initPathMap(p)
	pathMap["get"] = handleFunc
	return r
}

func (r *Router) Post(p string, handleFunc ...HandleFunc) IRouter {
	pathMap := r.initPathMap(p)
	pathMap["post"] = handleFunc
	return r
}

func (r *Router) MatchHandlersChain(method, p string) (HandlersChain, bool) {
	pathMap, pathOK := r.pathHandlerMap[p]
	if !pathOK {
		return nil, true
	}
	method = strings.ToLower(method)
	hc, hcOK := pathMap[method]
	return hc, hcOK
}
