package resp

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type StatusCode int

const (
	ok StatusCode = http.StatusOK
)

type listResponse struct {
	Code StatusCode `json:"code"`
	Msg  string     `json:"msg,omitempty"`
	Data datas      `json:"data"`
}

type datas struct {
	Total int64       `json:"total,omitempty"`
	Data  interface{} `json:"data"`
}

type dataResponse struct {
	Code StatusCode  `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

func Success(c echo.Context) error {
	var rjson struct {
		Code StatusCode `json:"code"`
		Msg  string     `json:"msg"`
	}

	rjson.Code = ok
	rjson.Msg = "success"

	return c.JSON(http.StatusOK, rjson)
}

func Msg(msg string, code StatusCode) error {
	return Error{
		Err:      errors.New(msg),
		HttpCode: http.StatusCreated,
		Code:     code,
		Context:  nil,
	}
}

func DBErr(err error, code StatusCode) error {
	return Error{
		Err:      err,
		HttpCode: http.StatusBadRequest,
		Code:     code,
		Context:  nil,
	}
}

func Unauthorized(err error) error {
	return Error{
		Err:      err,
		HttpCode: http.StatusUnauthorized,
		Code:     http.StatusUnauthorized,
		Context:  nil,
	}
}

func ParamErr(err error, code StatusCode) error {
	return Error{
		Err:      err,
		HttpCode: http.StatusBadRequest,
		Code:     code,
		Context:  nil,
	}
}

func ListResponse(arr interface{}, total int64, c echo.Context) error {
	if arr == nil {
		arr = make([]interface{}, 0)
	} else if reflect.ValueOf(arr).IsNil() {
		arr = make([]interface{}, 0)
	}
	r := listResponse{
		Data: datas{
			Data:  arr,
			Total: total,
		},
		Code: ok,
	}
	return c.JSON(http.StatusOK, r)
}

func Response(data interface{}, c echo.Context) error {
	r := dataResponse{
		Data: data,
		Code: ok,
	}

	return c.JSON(http.StatusOK, r)
}

func StatusResponse(status int, data interface{}, code StatusCode, c echo.Context) error {
	r := dataResponse{
		Data: data,
		Code: code,
	}
	return c.JSON(status, r)
}

type ErrorHandler func(ctx echo.Context)

func EchoErrorHandler(handlerFunc ...ErrorHandler) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		for _, v := range handlerFunc {
			v(c)
		}
		if err == nil {
			return
		}
		var rjson struct {
			Code StatusCode `json:"code"`
			Msg  string     `json:"msg"`
		}

		if terr, ok := err.(*Error); ok {
			rjson.Code = terr.Code
			rjson.Msg = terr.Err.Error()
			c.JSON(terr.HttpCode, rjson)

		} else {
			c.JSON(http.StatusBadRequest, err)
		}
	}
}

type Error struct {
	Err      error
	HttpCode int
	Code     StatusCode
	Context  echo.Context
}

func (e Error) Error() string {
	return e.Err.Error()
}
