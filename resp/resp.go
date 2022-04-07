package resp

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/5xxxx/log"

	"github.com/labstack/echo/v4"
	"github.com/marmotedu/errors"
)

type listResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data datas  `json:"data"`
}

type datas struct {
	Total int64       `json:"total,omitempty"`
	List  interface{} `json:"list"`
}

type dataResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data"`
}

func Success(c echo.Context) error {
	var rjson struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	rjson.Code = http.StatusOK
	rjson.Msg = "success"

	return c.JSON(http.StatusOK, rjson)
}

func Msg(m string) error {
	return errors.WithCode(ErrDatabase, m)
}

func DBErr(err error) error {
	return errors.WithCode(ErrDatabase, err.Error())
}

func Unauthorized(err error) error {
	return errors.WithCode(ErrPermissionDenied, err.Error())
}

func SignatureInvalid(err error) error {
	return errors.WithCode(ErrSignatureInvalid, err.Error())
}

func ParamErr(err error) error {
	return errors.WithCode(ErrValidation, err.Error())
}

func UnknownErr(err error) error {
	return errors.WithCode(ErrUnknown, err.Error())
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
		Code: http.StatusOK,
	}
	return c.JSON(http.StatusOK, r)
}

func Response(data interface{}, c echo.Context) error {
	r := dataResponse{
		Data: data,
		Msg:  "success",
		Code: http.StatusOK,
	}

	return c.JSON(http.StatusOK, r)
}

type ErrorHandler func(ctx echo.Context)

func EchoErrorHandler(log log.Logger, handlerFunc ...ErrorHandler) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		for _, v := range handlerFunc {
			v(c)
		}

		if err == nil {
			return
		}
		if log != nil {
			log.Error(fmt.Sprintf("%v", err))
		}

		var rjson struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}

		coder := errors.ParseCoder(err)
		if coder.Code() != 1 {
			rjson.Code = coder.Code()
			rjson.Msg = coder.String()
			c.JSON(coder.HTTPStatus(), rjson)
		} else {
			c.JSON(http.StatusBadRequest, err)
		}
	}
}
