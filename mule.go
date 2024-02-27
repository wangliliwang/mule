package mule

import "net/http"

func New() *App {
	return &App{NewRouter()}
}

// 实现 http.Handler
type App struct {
	IRouter
}

func (a *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		req:       request,
		resWriter: writer,
	}
	method := ctx.req.Method
	rPath := ctx.req.URL.Path
	handlersChain, ok := a.MatchHandlersChain(method, rPath)
	if !ok {
		ctx.SetResponse(http.StatusNotFound, "application/json", []byte{})
	}
	for _, handler := range handlersChain {
		handler(ctx)
	}
}

func (a *App) Run(addr string) error {
	return http.ListenAndServe(addr, a)
}

type HandleFunc func(ctx *Context)

type HandlersChain []HandleFunc
