package mule

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
)

type Context struct {
	request  *Request
	response *Response
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

//func (c *Context) Deadline() (deadline time.Time, ok bool) {
//	return c.request.Context().Deadline()
//}
//
//func (c *Context) Done() <-chan struct{} {
//	return c.request.Context().Done()
//}
//
//func (c *Context) Err() error {
//	return c.request.Context().Err()
//}
//
//func (c *Context) Value(key any) any {
//	return c.request.Context().Value(key)
//}

func (c *Context) SetResponse(httpCode int, contentType, content string) {
	c.response.SetStatusCode(httpCode)
	c.response.SetHeader("Content-Type", contentType)
	c.response.SetBody(content)
	c.response.Flush()
}

func (c *Context) parse(obj any, raw []byte) error {
	//TODO implement me
	unmarshalErr := json.Unmarshal(raw, obj)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	validateErr := validateStruct(obj)
	if validateErr != nil {
		return validateErr
	}
	return nil
}

// TODO(wangli) fix me 1. 支持传入多个值 2. 支持query tag
func (c *Context) ParseQuery(obj any) error {
	values := c.request.Query
	singleValues := make(map[string]string)
	for key, value := range values {
		if len(value) > 0 {
			singleValues[key] = value[0]
		}
	}
	bs, marshalErr := json.Marshal(singleValues)
	if marshalErr != nil {
		return marshalErr
	}
	return c.parse(obj, bs)
}

func (c *Context) ParseJSON(obj any) error {
	return c.parse(obj, []byte(c.request.Body))
}

var validate = validator.New()

func validateStruct(a any) error {
	validateErr := validate.Struct(a)
	if validateErr != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(validateErr, &invalidValidationError) {
			return fmt.Errorf("validate Error: %w", validateErr)
		}
		errorMessage := "param validate error: "
		for index, err := range validateErr.(validator.ValidationErrors) {
			errorMessage += fmt.Sprintf("%d - field [%s] validate failed by rule [%s]; ", index+1, err.StructNamespace(), err.ActualTag())
		}
		return errors.New(errorMessage)
	}
	return nil
}
