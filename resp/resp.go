package resp

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
)

type StatusCode int

const (
	ok           StatusCode = http.StatusOK
	db           StatusCode = 407
	msg          StatusCode = 201
	param        StatusCode = 405
	unauthorized StatusCode = 401
	unknown      StatusCode = 501
)

type listResponse struct {
	Code StatusCode `json:"code"`
	Msg  string     `json:"msg,omitempty"`
	Data datas      `json:"data"`
}

type datas struct {
	Total int64       `json:"total,omitempty"`
	List  interface{} `json:"list"`
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

func BizErrStatus(err error, code StatusCode) error {
	return Error{
		Err:      errors.New(fmt.Sprintf("%+v", err)),
		HttpCode: http.StatusBadRequest,
		Code:     code,
		Context:  nil,
	}
}

func Msg(m string) error {
	return Error{
		Err:      errors.New(m),
		HttpCode: http.StatusCreated,
		Code:     msg,
		Context:  nil,
	}
}

func DBErr(err error) error {
	return Error{
		Err:      errors.New(fmt.Sprintf("%+v", err)),
		HttpCode: http.StatusBadRequest,
		Code:     db,
		Context:  nil,
	}
}

func Unauthorized(err error) error {
	return Error{
		Err:      errors.New(fmt.Sprintf("%+v", err)),
		HttpCode: http.StatusUnauthorized,
		Code:     unauthorized,
		Context:  nil,
	}
}

func ParamErr(err error) error {
	return Error{
		Err:      errors.New(fmt.Sprintf("%+v", err)),
		HttpCode: http.StatusBadRequest,
		Code:     param,
		Context:  nil,
	}
}

func UnknownErr(err error) error {
	return Error{
		Err:      errors.New(fmt.Sprintf("%+v", err)),
		HttpCode: http.StatusBadRequest,
		Code:     unknown,
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
			List:  arr,
			Total: total,
		},
		Code: ok,
	}
	return c.JSON(http.StatusOK, r)
}

func Response(data interface{}, c echo.Context) error {
	r := dataResponse{
		Data: data,
		Msg: "success",
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

		if terr, ok := err.(Error); ok {
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
