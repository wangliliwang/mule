package mule

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"time"
)

type Context struct {
	req       *http.Request
	resWriter http.ResponseWriter
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.req.Context().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.req.Context().Done()
}

func (c *Context) Err() error {
	return c.req.Context().Err()
}

func (c *Context) Value(key any) any {
	return c.req.Context().Value(key)
}

func (c *Context) SetResponse(httpCode int, contentType string, content []byte) {
	c.resWriter.WriteHeader(httpCode)
	c.resWriter.Header().Set("Content-Type", contentType)
	c.resWriter.Write(content)
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
	values := c.req.URL.Query()
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
	if c.req.Body == nil {
		return errors.New("invalid request")
	}
	bd, readErr := io.ReadAll(c.req.Body)
	if readErr != nil {
		return readErr
	}
	return c.parse(obj, bd)
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
