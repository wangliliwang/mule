package mule

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

type handler1Req struct {
	Name string `json:"name" validate:"required"`
	Age  string `json:"age" validate:"required"`
}

func testHandler1(ctx *Context) {
	req := &handler1Req{}
	parseErr := ctx.ParseQuery(req)
	if parseErr != nil {
		ctx.SetResponse(http.StatusBadRequest, "application/json", []byte(parseErr.Error()))
		return
	}
	ctx.SetResponse(http.StatusOK, "application/json", []byte(fmt.Sprintf(`{"name": "%s"}`, req.Name)))
}

type handler2Req struct {
	Name string `json:"name" validate:"required"`
	Age  string `json:"age" validate:"required"`
}

func testHandler2(ctx *Context) {
	req := &handler2Req{}
	parseErr := ctx.ParseJSON(req)
	log.Printf("in handler2 %+v\n", req)
	if parseErr != nil {
		ctx.SetResponse(http.StatusBadRequest, "application/json", []byte(parseErr.Error()))
		return
	}
	ctx.SetResponse(http.StatusOK, "application/json", []byte(fmt.Sprintf(`{"name": "%s"}`, req.Name)))
}

func TestApp(t *testing.T) {
	app := New()
	app.Get("/hello", testHandler1)
	app.Post("/world", testHandler2)
	log.Fatal(app.Run(":8888"))
}
