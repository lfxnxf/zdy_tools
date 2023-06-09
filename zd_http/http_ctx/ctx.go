package http_ctx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lfxnxf/zdy_tools/trace"
	"github.com/lfxnxf/zdy_tools/zd_error"
	"net/http"
	"reflect"
	"strconv"
)

type HttpHandler func(ctx *HttpContext)

type WrapResp struct {
	Code      string      `json:"code"`
	Msg       string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestId string      `json:"requestId"`
}

const (
	RespKey = "response_data"
)

func NewWrapResp(data interface{}, err error, traceId string) WrapResp {
	var e = zd_error.Cause(err)
	return WrapResp{
		Code:      e.Code(),
		Msg:       e.Message(),
		Data:      data,
		RequestId: traceId,
	}
}

type HttpContext struct {
	*gin.Context
}

func (c *HttpContext) WriteJson(data interface{}, err error) {
	w := NewWrapResp(data, err, trace.ExtraTraceID(c))
	if w.Code != "" {
		c.JSON(http.StatusInternalServerError, w)
	} else {
		c.JSON(http.StatusOK, w)
	}
	//c.Set(RespKey, w)
	//c.JSON(http.StatusOK, w)
	//c.Abort()
}

func (c *HttpContext) WriteParamsError(err error, data interface{}) {
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
		c.WriteJson(nil, zd_error.AddSpecialError(zd_error.ParamsErrorCode, msg))
	} else {
		c.WriteJson(nil, zd_error.ParamsError)
	}
}

func (c *HttpContext) DefaultQueryInt64(key string, d int64) int64 {
	val := c.DefaultQuery(key, "0")
	valInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil || valInt == 0 {
		return d
	}
	return valInt
}
