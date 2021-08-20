package resp

import (
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
	Msg  string     `json:"msg"`
	Data datas      `json:"data"`
}

type datas struct {
	Total int64       `json:"total,omitempty"`
	Data  interface{} `json:"data"`
}

type dataResponse struct {
	Code StatusCode  `json:"code"`
	Msg  string      `json:"msg"`
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

func Msg(status int, msg string, code StatusCode, c echo.Context) error {
	var rjson struct {
		Code StatusCode `json:"code"`
		Msg  string     `json:"msg"`
	}
	rjson.Code = code
	rjson.Msg = msg
	return c.JSON(status, rjson)
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
