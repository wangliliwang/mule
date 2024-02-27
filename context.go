package mule

import (
	"net/http"
	"time"
)

type Context struct {
	req       *http.Request
	resWriter http.ResponseWriter
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Done() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Err() error {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Value(key any) any {
	//TODO implement me
	panic("implement me")
}

func (c *Context) SetResponse(httpCode int, contentType string, content []byte) {
	c.resWriter.WriteHeader(httpCode)
	c.resWriter.Header().Set("Content-Type", contentType)
	c.resWriter.Write(content)
}

func (c *Context) Bind(obj any) error {
	//TODO implement me
	panic("implement me")
}
