package zd_http

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/lfxnxf/zdy_tools/trace"
	"github.com/lfxnxf/zdy_tools/zd_error"
)

type WrapResp struct {
	Code      string      `json:"code"`
	Msg       string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestId string      `json:"requestId"`
}

func newWrapResp(data interface{}, err error, traceId string) WrapResp {
	var e = zd_error.Cause(err)
	return WrapResp{
		Code:      e.Code(),
		Msg:       e.Message(),
		Data:      data,
		RequestId: traceId,
	}
}

func WriteJson(c *gin.Context, data interface{}, err error) {
	w := newWrapResp(data, err, trace.ExtraTraceID(c))
	// 所有都返回200
	c.JSON(http.StatusOK, w)
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
