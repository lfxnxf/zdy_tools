package zd_http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lfxnxf/zdy_tools/zd_error"
	"net/http"
	"reflect"
)

type WrapResp struct {
	Code int         `json:"error"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newWrapResp(data interface{}, err error) WrapResp {
	var e = zd_error.Cause(err)
	return WrapResp{
		Code: e.Code(),
		Msg:  e.Message(),
		Data: data,
	}
}

func WriteJson(c *gin.Context, data interface{}, err error) {
	w := newWrapResp(data, err)
	if w.Code != 0 {
		c.JSON(http.StatusInternalServerError, w)
	} else {
		c.JSON(http.StatusOK, w)
	}
	//c.JSON(http.StatusOK, w)
	c.Abort()
}

func WriteParamsError(c *gin.Context, err error, data interface{}) {
	var isParamsError bool
	var msg string
	obj := reflect.TypeOf(data)
	if validErr, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validErr {
			if f, exist := obj.Elem().FieldByName(e.Field()); exist && len(f.Tag.Get("msg")) > 0 {
				isParamsError = true
				msg = fmt.Sprintf("参数错误，%s", f.Tag.Get("msg"))
			}
		}
	}
	if isParamsError {
		WriteJson(c, nil, zd_error.AddSpecialError(zd_error.ParamsErrorCode, msg))
	} else {
		WriteJson(c, nil, zd_error.ParamsError)
	}
}