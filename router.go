package mule

type IRouter interface {
	// Add
	Handle(string, string, ...HandleFunc) IRouter
	Get(string, ...HandleFunc) IRouter
	Post(string, ...HandleFunc) IRouter

	// Query
	MatchHandlersChain(string, string) (HandlersChain, bool)
}

type Router struct {
	pathHandlerMap map[string]HandlersChain
}

func NewRouter() IRouter {
	return &Router{pathHandlerMap: make(map[string]HandlersChain)}
}

func (r *Router) Handle(s string, s2 string, handleFunc ...HandleFunc) IRouter {
	//TODO implement me
	panic("implement me")
}

func (r *Router) Get(s string, handleFunc ...HandleFunc) IRouter {
	//TODO implement me
	panic("implement me")
}

func (r *Router) Post(s string, handleFunc ...HandleFunc) IRouter {
	//TODO implement me
	panic("implement me")
}

func (r *Router) MatchHandlersChain(method, p string) (HandlersChain, bool) {
	//TODO implement me
	panic("implement me")
}
